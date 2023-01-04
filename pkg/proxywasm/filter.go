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
	"emperror.dev/errors"
	"github.com/go-logr/logr"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

type filterContext struct {
	api.ABIContext

	id          int32
	plugin      api.WasmPlugin
	rootContext api.Context
	logger      logr.Logger
}

func NewFilterContext(plugin api.WasmPlugin, properties api.PropertyHolder) (api.FilterContext, error) {
	instance, err := plugin.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hostOptions := []HostFunctionsOption{
		SetHostFunctionsLogger(plugin.Logger()),
	}

	if val, ok := GetBaseContext().Get("metric.handler"); ok {
		if mh, ok := val.(api.MetricHandler); ok {
			hostOptions = append(hostOptions, SetHostFunctionsMetricHandler(mh))
		}
	}

	context := &filterContext{
		ABIContext: abi.NewContext(NewHostFunctions(
			NewPropertyHolderWrapper(properties, plugin.GetWasmInstanceContext(instance).GetProperties()),
			hostOptions...,
		), instance),
		logger:      plugin.Logger(),
		id:          plugin.Context().NewContextID(),
		rootContext: plugin.Context(),
		plugin:      plugin,
	}

	context.Lock()
	defer context.Unlock()

	if err := context.GetExports().ProxyOnContextCreate(context.ID(), context.rootContext.ID()); err != nil {
		return nil, errors.WrapIfWithDetails(err, "could not create context", "id", context.ID)
	}

	plugin.RegisterFilterContext(instance, context)

	return context, nil
}

func (c *filterContext) GetABIContext() api.ContextHandler {
	return c.ABIContext
}

func (c *filterContext) Logger() logr.Logger {
	return c.logger
}

func (c *filterContext) Lock() {
	c.ABIContext.GetInstance().Lock(c.ABIContext)
}

func (c *filterContext) Unlock() {
	c.ABIContext.GetInstance().Unlock()
}

func (c *filterContext) Close() {
	c.Unlock()
	c.plugin.ReleaseInstance(c.ABIContext.GetInstance())
	c.plugin.UnregisterFilterContext(c.ABIContext.GetInstance(), c)
}

func (c *filterContext) ID() int32 {
	return c.id
}

func (c *filterContext) RootContext() api.Context {
	return c.rootContext
}
