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
	"runtime"
	"sync"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/pborman/uuid"

	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

type wasmPluginManager struct {
	vms    api.VMStore
	logger logr.Logger

	plugins sync.Map
}

func NewWasmPluginManager(vms api.VMStore, logger logr.Logger) api.WasmPluginManager {
	return &wasmPluginManager{
		vms:    vms,
		logger: logger,

		plugins: sync.Map{},
	}
}

func (m *wasmPluginManager) Logger() logr.Logger {
	return m.logger
}

func (m *wasmPluginManager) Delete(config api.WasmPluginConfig) error {
	plugin, err := m.Get(config)
	if err != nil {
		return err
	}

	plugin.Close()

	m.plugins.Delete(config.Key())

	return nil
}

func (m *wasmPluginManager) Get(config api.WasmPluginConfig) (api.WasmPlugin, error) {
	if val, ok := m.plugins.Load(config.Key()); ok {
		if plugin, ok := val.(api.WasmPlugin); ok {
			return plugin, nil
		}
	}

	return nil, errors.NewWithDetails("plugin is not found", "name", config.Name)
}

func (m *wasmPluginManager) GetOrCreate(config api.WasmPluginConfig) (api.WasmPlugin, error) {
	if val, ok := m.plugins.Load(config.Key()); ok {
		if plugin, ok := val.(api.WasmPlugin); ok {
			m.Logger().V(3).Info("getOrCreate plugin", "cached", true, "key", config.Key())
			return plugin, nil
		}
	}

	m.Logger().V(3).Info("getOrCreate plugin", "cached", false, "key", config.Key())

	return m.Add(config)
}

func (m *wasmPluginManager) Add(config api.WasmPluginConfig) (api.WasmPlugin, error) {
	var err error

	if config.Name == "" {
		config.Name = uuid.New()
	}

	plugin := &wasmPlugin{
		config: config,
		logger: m.logger.WithName("wasmPlugin").WithName(config.Name),

		borrowedInstances: make(map[api.WasmInstance]uint32),

		filterContexts:   map[api.WasmInstance]*sync.Map{},
		instanceContexts: &sync.Map{},
	}

	if plugin.config.InstanceCount <= 0 {
		plugin.config.InstanceCount = uint32(runtime.NumCPU())
	}

	plugin.vm, err = m.vms.GetOrCreateVM(config.VMConfig)
	if err != nil {
		return nil, errors.WrapIf(err, "could not get vm")
	}

	plugin.vm.Acquire(plugin)

	wasmBytes, err := config.VMConfig.Code.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if module, err := plugin.vm.NewModule(wasmBytes); err != nil {
		return nil, errors.WithStack(err)
	} else {
		plugin.module = module
	}

	// plugin.module.SetLogger(plugin.logger.WithName("wasm"))

	plugin.ctx = GetBaseContext().GetOrCreateContext(plugin.config.RootID)

	plugin.ctx.Set("plugin_root_id", config.RootID)
	plugin.ctx.Set("plugin_vm_id", config.VMConfig.ID)
	plugin.ctx.Set("plugin_name", config.Name)

	PluginProperty(plugin.ctx).Set(plugin)

	if err := plugin.setPluginConfig(); err != nil {
		return nil, err
	}

	plugin.EnsureInstances(config.InstanceCount)

	m.plugins.Store(plugin.config.Key(), plugin)

	return plugin, nil
}
