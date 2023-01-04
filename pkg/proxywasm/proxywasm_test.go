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

package proxywasm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2/klogr"

	"github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero"
	"github.com/cisco-open/nasp/pkg/dotn"
	"github.com/cisco-open/nasp/pkg/proxywasm"
	"github.com/cisco-open/nasp/pkg/proxywasm/api"
	"github.com/cisco-open/nasp/pkg/proxywasm/testdata"
)

func TestProxyOnTickAndSetEffectiveContext(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	logger := klogr.New()

	runtimeCreators := proxywasm.NewRuntimeCreatorStore()
	runtimeCreators.Set("wazero", func() api.WasmRuntime {
		return wazero.NewVM(context.Background(), wazero.VMWithLogger(logger))
	})

	vms := proxywasm.NewVMStore(runtimeCreators, logger)
	pm := proxywasm.NewWasmPluginManager(vms, logger)

	// the ontick wasm plugin will run at every second and collect the 'test-id' property
	// from every filter context and eventually store it in the 'test-ids' property.
	filter := api.WasmPluginConfig{
		Name: "ontick",
		VMConfig: api.WasmVMConfig{
			Runtime: "wazero",
			Code:    proxywasm.NewFileDataSource(testdata.Filters, "ontick.wasm"),
		},
		InstanceCount: 1,
	}

	plugin, err := pm.GetOrCreate(filter)
	require.Nil(err)

	for i := 0; i < 5; i++ {
		properties := dotn.New()

		_, err = proxywasm.NewFilterContext(plugin, properties)
		require.Nil(err)

		properties.Set("test-id", fmt.Sprintf("context-%d", i))
	}

	inst, err := plugin.GetInstance()
	require.Nil(err)

	require.Eventually(func() bool {
		if v, ok := plugin.GetWasmInstanceContext(inst).GetProperties().Get("test-ids"); ok {
			if ids, ok := v.(string); ok {
				return "context-0,context-1,context-2,context-3,context-4" == ids
			}
		}

		return false
	}, time.Second*5, time.Millisecond*500, "'test-ids' property is not set properly")
}
