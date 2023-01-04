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
	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

type tickerDoneChannelProperty struct {
	api.PropertyHolder
	key string
}

func TickerDoneChannelProperty(p api.PropertyHolder) tickerDoneChannelProperty {
	return tickerDoneChannelProperty{p, "ticker.done_channel"}
}

func (p tickerDoneChannelProperty) Get() (chan bool, bool) {
	if v, ok := p.PropertyHolder.Get(p.key); ok {
		if done, ok := v.(chan bool); ok {
			done <- true
		}
	}

	return nil, false
}

func (p tickerDoneChannelProperty) Set(done chan bool) {
	p.PropertyHolder.Set(p.key, done)
}

type pluginProperty struct {
	api.PropertyHolder
	key string
}

func PluginProperty(p api.PropertyHolder) pluginProperty {
	return pluginProperty{p, "plugin"}
}

func (p pluginProperty) Get() (api.WasmPlugin, bool) {
	if v, ok := p.PropertyHolder.Get(p.key); ok {
		if plugin, ok := v.(api.WasmPlugin); ok {
			return plugin, true
		}
	}

	return nil, false
}

func (p pluginProperty) Set(plugin api.WasmPlugin) {
	p.PropertyHolder.Set(p.key, plugin)
}

type rootABIContextProperty struct {
	api.PropertyHolder
	key string
}

func RootABIContextProperty(p api.PropertyHolder) rootABIContextProperty {
	return rootABIContextProperty{p, "root_abi_context"}
}

func (p rootABIContextProperty) Get() (api.ContextHandler, bool) {
	if v, ok := p.PropertyHolder.Get(p.key); ok {
		if ctx, ok := v.(api.ContextHandler); ok {
			return ctx, true
		}
	}

	return nil, false
}

func (p rootABIContextProperty) Set(ctx api.ContextHandler) {
	p.PropertyHolder.Set(p.key, ctx)
}
