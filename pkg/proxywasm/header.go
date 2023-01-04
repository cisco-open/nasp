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
	"strings"

	"github.com/go-logr/logr"

	pwapi "github.com/banzaicloud/proxy-wasm-go-host/api"

	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

// convert HeaderMap to api.HeaderMap.
type Headers struct {
	headers api.HeaderMap
	logger  logr.Logger
}

type HeaderMap interface {
	pwapi.HeaderMap

	Flatten() map[string]string
}

func NewHeaders(headers api.HeaderMap, logger logr.Logger) HeaderMap {
	return &Headers{
		headers: headers,
		logger:  logger,
	}
}

func (h *Headers) Get(key string) (value string, found bool) {
	value, found = h.headers.Get(key)
	if !found {
		return
	}

	h.logger.V(2).Info("get header", "key", key, "value", value)

	return
}

func (h *Headers) Set(key, value string) {
	h.logger.V(2).Info("set header", "key", key, "value", value)

	h.headers.Set(key, value)
}

func (h *Headers) Add(key, value string) {
	h.logger.V(2).Info("add header", "key", key, "value", value)

	h.headers.Add(key, value)
}

func (h *Headers) Del(key string) {
	h.logger.V(2).Info("del header", "key", key)

	h.headers.Del(key)
}

func (h *Headers) Range(f func(key string, value string) bool) {
	for _, k := range h.headers.Keys() {
		if v, found := h.headers.Get(k); found {
			f(k, v)
		}
	}
}

func (h *Headers) Clone() pwapi.HeaderMap {
	return &Headers{
		headers: h.headers.Clone(),
	}
}

func (h *Headers) ByteSize() uint64 {
	var size int

	for key, values := range h.headers.Raw() {
		size += len(key)
		for _, value := range values {
			size += len(value)
		}
	}

	return uint64(size)
}

func (h *Headers) Flatten() map[string]string {
	flatten := make(map[string]string)

	for key, values := range h.headers.Raw() {
		flatten[key] = strings.Join(values, ",")
	}

	return flatten
}
