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

package api

import (
	"github.com/go-logr/logr"
	v1 "mosn.io/proxy-wasm-go-host/proxywasm/v1"
)

type FilterContext interface {
	v1.ContextHandler

	Logger() logr.Logger
	Lock()
	Unlock()
	Close()
	ID() int32
	RootContext() Context
	GetABIContext() v1.ContextHandler
}
