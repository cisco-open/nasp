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
	"fmt"
	"sync"

	"emperror.dev/errors"
	"github.com/go-logr/logr"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"

	"github.com/cisco-open/nasp/pkg/dotn"
	"github.com/cisco-open/nasp/pkg/proxywasm/api"
)

type wasmPlugin struct {
	config api.WasmPluginConfig
	logger logr.Logger

	module api.WasmModule
	vm     api.WasmVM
	ctx    api.Context

	instanceCount  uint32
	instances      []api.WasmInstance
	instancesIndex uint32

	borrowedInstancesCount uint32
	borrowedInstances      map[api.WasmInstance]uint32

	lock sync.RWMutex

	filterContexts     map[api.WasmInstance]*sync.Map
	filterContextsLock sync.RWMutex
	instanceContexts   *sync.Map
}

func (p *wasmPlugin) InstanceCount() uint32 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.instanceCount
}

func (p *wasmPlugin) Logger() logr.Logger {
	return p.logger
}

func (p *wasmPlugin) Report() {
	p.logger.V(2).Info("current index", "index", p.instancesIndex)
	p.logger.V(2).Info("borrowed instances count", "count", p.borrowedInstancesCount)
}

func (p *wasmPlugin) Name() string {
	return p.config.Name
}

func (p *wasmPlugin) ID() string {
	return fmt.Sprintf("%s-%s", p.config.RootID, p.Name())
}

func (p *wasmPlugin) GetInstance() (api.WasmInstance, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for i := 0; i < len(p.instances); i++ {
		idx := p.instancesIndex % uint32(len(p.instances))
		p.instancesIndex++

		instance := p.instances[idx]
		if !instance.Acquire() {
			continue
		}

		p.borrowedInstancesCount++
		p.borrowedInstances[instance]++

		return instance, nil
	}

	return nil, errors.New("could not get instance")
}

func (p *wasmPlugin) ReleaseInstance(instance api.WasmInstance) {
	p.lock.Lock()
	defer p.lock.Unlock()

	instance.Release()

	if _, ok := p.borrowedInstances[instance]; ok {
		p.borrowedInstances[instance]--
	}

	p.borrowedInstancesCount--
	p.instancesIndex--
}

func (p *wasmPlugin) Exec(f func(instance api.WasmInstance) bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for _, instance := range p.instances {
		if !f(instance) {
			break
		}
	}
}

func (p *wasmPlugin) Close() {
	p.lock.Lock()

	// release every borrowed instance
	for instance, bi := range p.borrowedInstances {
		for i := bi; i > 0; i-- {
			instance.Release()
			p.borrowedInstancesCount--
			p.instancesIndex--
		}
		delete(p.borrowedInstances, instance)
	}

	p.lock.Unlock()

	p.EnsureInstances(0)

	p.vm.Release(p)
}

func (p *wasmPlugin) VM() api.WasmVM {
	return p.vm
}

func (p *wasmPlugin) EnsureInstances(desired uint32) uint32 {
	p.config.InstanceCount = desired

	actual := p.InstanceCount()

	p.logger.V(2).Info("ensure instances", "actual", actual, "desired", desired)

	if actual == desired {
		return actual
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	if desired < actual {
		for i := int(desired); i < len(p.instances); i++ {
			p.logger.V(2).Info("drop instance")
			if err := p.stopInstance(p.instances[i]); err != nil {
				p.logger.Error(err, "could not stop instance")
			}
			p.instances[i].Stop()
			p.instances[i] = nil
		}

		p.instances = p.instances[:desired]
		p.instanceCount = desired

		return p.instanceCount
	}

	newInstance := make([]api.WasmInstance, 0)
	numToCreate := desired - actual

	for i := 0; i < int(numToCreate); i++ {
		p.logger.V(2).Info("create new instance")

		instance, err := p.module.NewInstance()
		if err != nil {
			p.logger.Error(err, "could not create new instance")
			continue
		}

		if err := instance.Start(); err != nil {
			p.logger.Error(err, "could not start instance")
			continue
		}

		if err := p.startInstance(instance); err != nil {
			p.logger.Error(err, "could not start instance")
			continue
		}

		newInstance = append(newInstance, instance)
	}

	p.instances = append(p.instances, newInstance...)
	p.instanceCount += uint32(len(newInstance))

	return p.instanceCount
}

func (p *wasmPlugin) setPluginConfig() error {
	p.ctx.Set("plugin_config.bytes", "")
	p.ctx.Set("plugin_config.size", 0)

	if p.config.Configuration == nil {
		return nil
	}

	json, err := p.config.Configuration.Marshal()
	if err != nil {
		return errors.WrapIf(err, "could not marshal plugin config")
	}

	p.ctx.Set("plugin_config.bytes", json)
	p.ctx.Set("plugin_config.size", len(json))

	return nil
}

func (p *wasmPlugin) Context() api.Context {
	return p.ctx
}

func (p *wasmPlugin) getPluginConfig() (content []byte, size int) {
	if val, ok := p.ctx.Get("plugin_config.bytes"); ok {
		if b, ok := val.([]byte); ok {
			content = b
		}
	}

	if val, ok := p.ctx.Get("plugin_config.size"); ok {
		if s, ok := val.(int); ok {
			size = s
		}
	}

	return
}

func (p *wasmPlugin) GetWasmInstanceContext(instance api.WasmInstance) api.WasmInstanceContext {
	if v, ok := p.instanceContexts.Load(instance); ok {
		if iwc, ok := v.(api.WasmInstanceContext); ok {
			return iwc
		}
	}

	hostOptions := []HostFunctionsOption{
		SetHostFunctionsLogger(p.logger),
	}

	if val, ok := GetBaseContext().Get("metric.handler"); ok {
		if mh, ok := val.(api.MetricHandler); ok {
			hostOptions = append(hostOptions, SetHostFunctionsMetricHandler(mh))
		}
	}

	instanceProperties := NewPropertyHolderWrapper(dotn.New(), p.ctx)

	ctx := abi.NewContext(NewHostFunctions(instanceProperties, hostOptions...), instance)

	RootABIContextProperty(instanceProperties).Set(ctx)

	iwc := &instanceContext{
		instance:       instance,
		properties:     instanceProperties,
		contextHandler: ctx,
	}

	p.instanceContexts.Store(instance, iwc)

	return iwc
}

func (p *wasmPlugin) stopInstance(instance api.WasmInstance) error {
	ctx := p.GetWasmInstanceContext(instance).GetABIContext()

	instance.Lock(ctx)
	defer instance.Unlock()

	if err := StopWasmContext(p.ctx.ID(), ctx, p.logger); err != nil {
		return err
	}

	p.removeInstanceFromfilterContexts(instance)
	p.instanceContexts.Delete(instance)

	return nil
}

func (p *wasmPlugin) startInstance(instance api.WasmInstance) error {
	ctx := p.GetWasmInstanceContext(instance).GetABIContext()

	instance.Lock(ctx)
	defer instance.Unlock()

	if err := ctx.GetExports().ProxyOnContextCreate(p.ctx.ID(), 0); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnContextCreate", "contextID", p.ctx.ID())
	} else {
		p.logger.V(3).Info("root context is created successfully")
	}

	if res, err := ctx.GetExports().ProxyOnVmStart(p.ctx.ID(), 0); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnVmStart", "contextID", p.ctx.ID())
	} else if !res {
		return errors.NewWithDetails("unknown error at ProxyOnVmStart", "contextID", p.ctx.ID())
	} else {
		p.logger.V(3).Info("ProxyOnVmStart has run successfully", "contextID", p.ctx.ID())
	}

	_, size := p.getPluginConfig()
	if res, err := ctx.GetExports().ProxyOnConfigure(p.ctx.ID(), int32(size)); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnConfigure", "contextID", p.ctx.ID())
	} else if !res {
		return errors.NewWithDetails("unknown error at ProxyOnConfigure", "contextID", p.ctx.ID())
	} else {
		p.logger.V(3).Info("ProxyOnConfigure has run successfully")
	}

	return nil
}

func (p *wasmPlugin) RegisterFilterContext(instance api.WasmInstance, filterContext api.FilterContext) {
	p.filterContextsLock.Lock()
	defer p.filterContextsLock.Unlock()

	if _, ok := p.filterContexts[instance]; !ok {
		p.filterContexts[instance] = &sync.Map{}
	}

	p.logger.V(2).Info("register filter context", "id", filterContext.ID())

	p.filterContexts[instance].Store(filterContext.ID(), filterContext)
}

func (p *wasmPlugin) UnregisterFilterContext(instance api.WasmInstance, filterContext api.FilterContext) {
	p.filterContextsLock.Lock()
	defer p.filterContextsLock.Unlock()

	if _, ok := p.filterContexts[instance]; !ok {
		return
	}

	p.logger.V(2).Info("unregister filter context", "id", filterContext.ID())

	p.filterContexts[instance].Delete(filterContext.ID())
}

func (p *wasmPlugin) GetFilterContext(instance api.WasmInstance, id int32) (api.FilterContext, bool) {
	p.filterContextsLock.RLock()
	defer p.filterContextsLock.RUnlock()

	if _, ok := p.filterContexts[instance]; !ok {
		return nil, false
	}

	if v, ok := p.filterContexts[instance].Load(id); !ok {
		return nil, false
	} else if fc, ok := v.(api.FilterContext); ok {
		return fc, true
	}

	return nil, false
}

func (p *wasmPlugin) removeInstanceFromfilterContexts(instance api.WasmInstance) {
	p.filterContextsLock.Lock()
	defer p.filterContextsLock.Unlock()

	delete(p.filterContexts, instance)
}

type instanceContext struct {
	instance       api.WasmInstance
	properties     api.PropertyHolder
	contextHandler api.ContextHandler
}

func (iwc *instanceContext) GetInstance() api.WasmInstance {
	return iwc.instance
}

func (iwc *instanceContext) GetProperties() api.PropertyHolder {
	return iwc.properties
}

func (iwc *instanceContext) GetABIContext() api.ContextHandler {
	return iwc.contextHandler
}
