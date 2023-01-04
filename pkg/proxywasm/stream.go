// Copyright (c) 2022 Cisco and/or its affiliates. All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package proxywasm

import (
	"encoding/binary"
	"io"
	"net"

	"emperror.dev/errors"
	"github.com/go-logr/logr"

	pwapi "github.com/banzaicloud/proxy-wasm-go-host/api"

	"github.com/cisco-open/nasp/pkg/dotn"
	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

type streamHandler struct {
	pm api.WasmPluginManager

	filters []api.WasmPluginConfig
	plugins []api.WasmPlugin
}

func NewStreamHandler(pm api.WasmPluginManager, filters []api.WasmPluginConfig) (api.StreamHandler, error) {
	plugins := make([]api.WasmPlugin, 0)

	for _, filter := range filters {
		plugin, err := pm.GetOrCreate(filter)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, plugin)
	}

	return &streamHandler{
		pm: pm,

		filters: filters,
		plugins: plugins,
	}, nil
}

func (h *streamHandler) Logger() logr.Logger {
	return h.pm.Logger()
}

type stream struct {
	api.PropertyHolder

	direction api.ListenerDirection
	logger    logr.Logger

	filterNames map[string]struct{}

	filterContexts []api.FilterContext
}

func UsedFiltersStreamOption(filterNames ...string) api.StreamOption {
	return func(instance interface{}) {
		if stream, ok := instance.(*stream); ok {
			for _, name := range filterNames {
				stream.filterNames[name] = struct{}{}
			}
		}
	}
}

func (h *streamHandler) NewStream(direction api.ListenerDirection, options ...api.StreamOption) (api.Stream, error) {
	stream := &stream{
		PropertyHolder: NewPropertyHolderWrapper(dotn.New(), GetBaseContext()),

		direction: direction,
		logger:    h.pm.Logger(),

		filterNames: map[string]struct{}{},
	}

	for _, option := range options {
		option(stream)
	}

	stream.Set("listener_direction", func() string {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(direction))
		return string(b)
	}())

	stream.Set("listener_metadata", stream)

	contexts := []api.FilterContext{}

	filtersFound := map[string]struct{}{}

	for _, plugin := range h.plugins {
		if _, ok := stream.filterNames[plugin.Name()]; len(stream.filterNames) > 0 && !ok {
			continue
		}
		context, err := NewFilterContext(plugin, stream)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		contexts = append(contexts, context)
		filtersFound[plugin.Name()] = struct{}{}
	}

	stream.filterContexts = contexts

	return stream, nil
}

func (s *stream) HandleHTTPRequest(req api.HTTPRequest) error {
	s.Set("http.request", req)

	for _, filterContext := range s.filterContexts {
		action, err := func() (pwapi.Action, error) {
			filterContext.Lock()
			defer func() {
				filterContext.Unlock()
			}()

			return filterContext.GetExports().ProxyOnRequestHeaders(filterContext.ID(), int32(req.Header().Len()), 1)
		}()
		if err != nil {
			return err
		}
		if action == pwapi.ActionPause {
			return nil
		}
	}

	if req.Body() == nil {
		return nil
	}

	originalBody := req.Body()
	body := NewBuffer([]byte{})
	s.Set("request.body", body)
	req.SetBody(io.NopCloser(body))

	buf := make([]byte, 4096)
	var endOfStream int32
	for {
		nr, err := originalBody.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) {
			endOfStream = 1
		}
		if nr > 0 || endOfStream == 1 {
			body.Write(buf[:nr])
			for _, filterContext := range s.filterContexts {
				_, err := func() (pwapi.Action, error) {
					filterContext.Lock()
					defer func() {
						filterContext.Unlock()
					}()

					return filterContext.GetExports().ProxyOnRequestBody(filterContext.ID(), int32(nr), endOfStream)
				}()
				if err != nil {
					return err
				}
			}
		}

		if endOfStream == 1 {
			break
		}
	}

	return nil
}

func (s *stream) HandleHTTPResponse(resp api.HTTPResponse) error {
	s.Set("http.response", resp)

	for _, filterContext := range s.filterContexts {
		_, err := func() (pwapi.Action, error) {
			filterContext.Lock()
			defer func() {
				filterContext.Unlock()
			}()

			return filterContext.GetExports().ProxyOnResponseHeaders(filterContext.ID(), int32(resp.Header().Len()), 1)
		}()
		if err != nil {
			return err
		}
	}

	if resp.Body() == nil {
		return nil
	}

	originalBody := resp.Body()
	body := NewBuffer([]byte{})
	s.Set("response.body", body)
	resp.SetBody(io.NopCloser(body))

	buf := make([]byte, 4096)
	var endOfStream int32
	for {
		nr, err := originalBody.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) {
			endOfStream = 1
		}
		if nr > 0 || endOfStream == 1 {
			body.Write(buf[:nr])
			for _, filterContext := range s.filterContexts {
				_, err := func() (pwapi.Action, error) {
					filterContext.Lock()
					defer func() {
						filterContext.Unlock()
					}()

					return filterContext.GetExports().ProxyOnResponseBody(filterContext.ID(), int32(nr), endOfStream)
				}()
				if err != nil {
					return err
				}
			}
		}

		if endOfStream == 1 {
			break
		}
	}

	return nil
}

func (s *stream) HandleDownstreamData(conn net.Conn, n int, buf []byte) error {
	for _, filterContext := range s.filterContexts {
		_, err := func() (pwapi.Action, error) {
			filterContext.Lock()
			defer filterContext.Close()

			return filterContext.GetExports().ProxyOnDownstreamData(filterContext.ID(), int32(n), 0)
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *stream) HandleUpstreamData(conn net.Conn, n int, buf []byte) error {
	for _, filterContext := range s.filterContexts {
		_, err := func() (pwapi.Action, error) {
			filterContext.Lock()
			defer filterContext.Close()

			return filterContext.GetExports().ProxyOnUpstreamData(filterContext.ID(), int32(n), 0)
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *stream) HandleTCPNewConnection(conn net.Conn) error {
	for _, filterContext := range s.filterContexts {
		_, err := func() (pwapi.Action, error) {
			filterContext.Lock()
			defer filterContext.Close()

			return filterContext.GetExports().ProxyOnNewConnection(filterContext.ID())
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *stream) HandleTCPCloseConnection(conn net.Conn) error {
	for _, filterContext := range s.filterContexts {
		err := func() error {
			filterContext.Lock()
			defer filterContext.Close()

			if err := filterContext.GetExports().ProxyOnUpstreamConnectionClose(filterContext.ID(), int32(pwapi.PeerTypeRemote)); err != nil {
				return err
			}

			if err := filterContext.GetExports().ProxyOnDownstreamConnectionClose(filterContext.ID(), int32(pwapi.PeerTypeLocal)); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *stream) Direction() api.ListenerDirection {
	return s.direction
}

func (s *stream) Logger() logr.Logger {
	return s.logger
}

func (s *stream) Close() error {
	for _, filterContext := range s.filterContexts {
		err := func() error {
			filterContext.Lock()
			defer filterContext.Close()

			return StopWasmContext(filterContext.ID(), filterContext, filterContext.Logger())
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
