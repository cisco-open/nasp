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
	"sort"

	"emperror.dev/errors"
	"github.com/go-logr/logr"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

func StopWasmContext(contextID int32, abiContext api.ContextHandler, logger logr.Logger) error {
	if res, err := abiContext.GetExports().ProxyOnDone(contextID); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnDone", "contextID", contextID)
	} else if !res {
		return errors.NewWithDetails("unknown error at ProxyOnDone", "contextID", contextID)
	} else {
		logger.V(3).Info("ProxyOnDone has run successfully", "contextID", contextID)
	}

	if res := abiContext.GetImports().Done(); res != api.WasmResultOk {
		return errors.NewWithDetails("unknown error at Done", "contextID", contextID)
	} else {
		logger.V(3).Info("Done has run successfully", "contextID", contextID)
	}

	if err := abiContext.GetExports().ProxyOnLog(contextID); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnLog", "contextID", contextID)
	} else {
		logger.V(3).Info("ProxyOnLog has run successfully", "contextID", contextID)
	}

	if err := abiContext.GetExports().ProxyOnDelete(contextID); err != nil {
		return errors.WrapIfWithDetails(err, "error at ProxyOnDelete", "contextID", contextID)
	} else {
		logger.V(3).Info("ProxyOnDelete has run successfully", "contextID", contextID)
	}

	return nil
}

func serializeMapToPairs(pairs map[string]string) string {
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := make([]byte, 4+len(pairs)*8)
	// count
	binary.LittleEndian.PutUint32(buf[0:], uint32(len(pairs))) // 4
	i := 0
	for _, k := range keys {
		v := pairs[k]
		binary.LittleEndian.PutUint32(buf[4+(i*8):], uint32(len(k)))
		binary.LittleEndian.PutUint32(buf[8+(i*8):], uint32(len(v)))
		i++
	}
	for _, k := range keys {
		v := pairs[k]
		buf = append(buf, []byte(k+"\x00")...) //nolint:makezero
		buf = append(buf, []byte(v+"\x00")...) //nolint:makezero
	}

	return string(buf)
}
