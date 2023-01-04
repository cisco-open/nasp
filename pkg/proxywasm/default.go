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
	"time"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type DefaultHostFunctions struct {
	*abi.DefaultImportsHandler
}

func NewDefaultHostFunctions() api.ImportsHandler {
	return &DefaultHostFunctions{
		DefaultImportsHandler: &abi.DefaultImportsHandler{},
	}
}

func (d *DefaultHostFunctions) Wait() api.Action {
	return api.ActionContinue
}

// utils
func (d *DefaultHostFunctions) GetRootContextID() int32 {
	panic("GetRootContextID not implemented")
}

func (d *DefaultHostFunctions) GetVmConfig() api.IoBuffer {
	panic("GetVmConfig not implemented")
}

func (d *DefaultHostFunctions) GetPluginConfig() api.IoBuffer {
	panic("GetPluginConfig not implemented")
}

func (d *DefaultHostFunctions) Log(level api.LogLevel, msg string) api.WasmResult {
	panic("Log not implemented")
}

func (d *DefaultHostFunctions) SetEffectiveContextID(contextID int32) api.WasmResult {
	panic("SetEffectiveContextID not implemented")
}

func (d *DefaultHostFunctions) SetTickPeriodMilliseconds(tickPeriodMilliseconds int32) api.WasmResult {
	panic("SetTickPeriodMilliseconds not implemented")
}

func (d *DefaultHostFunctions) GetCurrentTimeNanoseconds() (int32, api.WasmResult) {
	nano := time.Now().Nanosecond()

	return int32(nano), api.WasmResultOk
}

func (d *DefaultHostFunctions) Done() api.WasmResult {
	panic("Done not implemented")
}

// l4

func (d *DefaultHostFunctions) GetDownStreamData() api.IoBuffer {
	panic("GetDownStreamData not implemented")
}

func (d *DefaultHostFunctions) GetUpstreamData() api.IoBuffer {
	panic("GetUpstreamData not implemented")
}

func (d *DefaultHostFunctions) ResumeDownstream() api.WasmResult {
	panic("ResumeDownstream not implemented")
}

func (d *DefaultHostFunctions) ResumeUpstream() api.WasmResult {
	panic("ResumeUpstream not implemented")
}

// http

func (d *DefaultHostFunctions) GetHttpRequestHeader() api.HeaderMap {
	panic("GetHttpRequestHeader not implemented")
}

func (d *DefaultHostFunctions) GetHttpRequestBody() api.IoBuffer {
	panic("GetHttpRequestBody not implemented")
}

func (d *DefaultHostFunctions) GetHttpRequestTrailer() api.HeaderMap {
	panic("GetHttpRequestTrailer not implemented")
}

func (d *DefaultHostFunctions) GetHttpResponseHeader() api.HeaderMap {
	panic("GetHttpResponseHeader not implemented")
}

func (d *DefaultHostFunctions) GetHttpResponseBody() api.IoBuffer {
	panic("GetHttpResponseBody not implemented")
}

func (d *DefaultHostFunctions) GetHttpResponseTrailer() api.HeaderMap {
	panic("GetHttpResponseTrailer not implemented")
}

func (d *DefaultHostFunctions) GetHttpCallResponseHeaders() api.HeaderMap {
	panic("GetHttpCallResponseHeaders not implemented")
}

func (d *DefaultHostFunctions) GetHttpCallResponseBody() api.IoBuffer {
	panic("GetHttpCallResponseBody not implemented")
}

func (d *DefaultHostFunctions) GetHttpCallResponseTrailer() api.HeaderMap {
	panic("GetHttpCallResponseTrailer not implemented")
}

func (d *DefaultHostFunctions) HttpCall(url string, headers api.HeaderMap, body api.IoBuffer, trailer api.HeaderMap, timeoutMilliseconds int32) (int32, api.WasmResult) {
	panic("HttpCall not implemented")
}

func (d *DefaultHostFunctions) ResumeHttpRequest() api.WasmResult {
	panic("ResumeHttpRequest not implemented")
}

func (d *DefaultHostFunctions) ResumeHttpResponse() api.WasmResult {
	panic("ResumeHttpResponse not implemented")
}

func (d *DefaultHostFunctions) SendHttpResp(respCode int32, respCodeDetail api.IoBuffer, respBody api.IoBuffer, additionalHeaderMap api.HeaderMap, grpcCode int32) api.WasmResult {
	panic("SendHttpResp not implemented")
}

// grpc

func (d *DefaultHostFunctions) OpenGrpcStream(grpcService string, serviceName string, method string) (int32, api.WasmResult) {
	panic("OpenGrpcStream not implemented")
}

func (d *DefaultHostFunctions) SendGrpcCallMsg(token int32, data api.IoBuffer, endOfStream int32) api.WasmResult {
	panic("SendGrpcCallMsg not implemented")
}

func (d *DefaultHostFunctions) CancelGrpcCall(token int32) api.WasmResult {
	panic("CancelGrpcCall not implemented")
}

func (d *DefaultHostFunctions) CloseGrpcCall(token int32) api.WasmResult {
	panic("CloseGrpcCall not implemented")
}

func (d *DefaultHostFunctions) GrpcCall(grpcService string, serviceName string, method string, data api.IoBuffer, timeoutMilliseconds int32) (int32, api.WasmResult) {
	panic("GrpcCall not implemented")
}

func (d *DefaultHostFunctions) GetGrpcReceiveInitialMetaData() api.HeaderMap {
	panic("GetGrpcReceiveInitialMetaData not implemented")
}

func (d *DefaultHostFunctions) GetGrpcReceiveBuffer() api.IoBuffer {
	panic("GetGrpcReceiveBuffer not implemented")
}

func (d *DefaultHostFunctions) GetGrpcReceiveTrailerMetaData() api.HeaderMap {
	panic("GetGrpcReceiveTrailerMetaData not implemented")
}

// foreign

func (d *DefaultHostFunctions) CallForeignFunction(funcName string, param []byte) ([]byte, api.WasmResult) {
	panic("CallForeignFunction not implemented")
}

func (d *DefaultHostFunctions) GetFuncCallData() api.IoBuffer {
	panic("GetFuncCallData not implemented")
}

// property

func (d *DefaultHostFunctions) GetProperty(key string) (string, api.WasmResult) {
	panic("GetProperty not implemented")
}

func (d *DefaultHostFunctions) SetProperty(key string, value string) api.WasmResult {
	panic("SetProperty not implemented")
}

// metric

func (d *DefaultHostFunctions) DefineMetric(metricType api.MetricType, name string) (int32, api.WasmResult) {
	panic("DefineMetric not implemented")
}

func (d *DefaultHostFunctions) IncrementMetric(metricID int32, offset int64) api.WasmResult {
	panic("IncrementMetric not implemented")
}

func (d *DefaultHostFunctions) RecordMetric(metricID int32, value int64) api.WasmResult {
	panic("RecordMetric not implemented")
}

func (d *DefaultHostFunctions) GetMetric(metricID int32) (int64, api.WasmResult) {
	panic("GetMetric not implemented")
}

func (d *DefaultHostFunctions) RemoveMetric(metricID int32) api.WasmResult {
	panic("RemoveMetric not implemented")
}

// custom extension

func (d *DefaultHostFunctions) GetCustomBuffer(bufferType api.BufferType) api.IoBuffer {
	panic("GetCustomBuffer not implemented")
}

func (d *DefaultHostFunctions) GetCustomHeader(mapType api.MapType) api.HeaderMap {
	panic("GetCustomHeader not implemented")
}
