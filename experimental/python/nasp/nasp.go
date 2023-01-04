//  Copyright (c) 2022 Cisco and/or its affiliates. All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//        https://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package main

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"k8s.io/klog/v2"

	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"unsafe"

	"emperror.dev/errors"

	"github.com/cisco-open/nasp/pkg/istio"
)

/*
#include <stdlib.h>

struct GoError {
	char *error_msg;
};

inline struct GoError *create_go_error(char *error_msg) {
	struct GoError *p;

	p = malloc(sizeof(struct GoError));
	p->error_msg = error_msg;

	return p;
};

struct GoHTTPHeaderItem {
	char *key;
	char *value;
};

struct GoHTTPHeaders {
	struct GoHTTPHeaderItem *items;
	unsigned int size;
};

struct GoHTTPResponse {
	int status_code;
	char *version;
	struct GoHTTPHeaders headers;
    char *body;
};

inline struct GoHTTPResponse *create_go_http_response(int status_code, char *version, struct GoHTTPHeaders headers, char *body) {
	struct GoHTTPResponse *p;

	p = malloc(sizeof(struct GoHTTPResponse));
	p->status_code = status_code;
	p->version = version;
	p->headers = headers;
	p->body = body;

	return p;
};

inline struct GoHTTPHeaderItem *new_http_headers(unsigned int count) {
	struct GoHTTPHeaderItem *p;

	p = malloc(sizeof(struct GoHTTPHeaderItem) * count);

	return p;
};


inline void set_http_header_at(struct GoHTTPHeaderItem *headers, struct GoHTTPHeaderItem header, unsigned int index) {
	headers[index] = header;
};

inline void free_http_headers(struct GoHTTPHeaderItem *headers, unsigned int size) {
	if(headers == NULL) {
		return;
	}

	for(int i = 0; i < size; i++) {
		free(headers[i].key);
		free(headers[i].value);
	}

	free(headers);
};

*/
import "C"

var logger = klog.NewKlogr()

var httpTransportPool = httpTransportPoolManager{}

//export NewHTTPTransport
func NewHTTPTransport(heimdallURLPtr, clientIDPtr, clientSecretPtr *C.char,
	usePushGateWay C._Bool, pushGatewayAddressPtr *C.char) (C.ulonglong, *C.struct_GoError) {
	heimdallURL := C.GoString(heimdallURLPtr)
	clientID := C.GoString(clientIDPtr)
	clientSecret := C.GoString(clientSecretPtr)

	hash := fnv.New64()
	hash.Write([]byte(heimdallURL))
	hash.Write([]byte(clientID))
	hash.Write([]byte(clientSecret))
	httpTransportIDHex := fmt.Sprintf("%x", hash.Sum(nil))
	httpTransportID, _ := strconv.ParseUint(httpTransportIDHex, 16, 64)

	if httpTransport, ok := httpTransportPool.Get(httpTransportID); ok {
		httpTransport.refCount++

		return C.ulonglong(httpTransportID), nil
	}

	var pgwConfig *istio.PushgatewayConfig
	if bool(usePushGateWay) {
		pgwConfig = &istio.PushgatewayConfig{
			Address:          C.GoString(pushGatewayAddressPtr),
			UseUniqueIDLabel: true,
		}
	}
	istioHandlerConfig := &istio.IstioIntegrationHandlerConfig{
		UseTLS:              true,
		IstioCAConfigGetter: istio.IstioCAConfigGetterHeimdall(heimdallURL, clientID, clientSecret, "v1"),
		PushgatewayConfig:   pgwConfig,
	}

	iih, err := istio.NewIstioIntegrationHandler(istioHandlerConfig, logger)

	if err != nil {
		return C.ulonglong(0), cGoError(err)
	}

	transport, err := iih.GetHTTPTransport(http.DefaultTransport)
	if err != nil {
		return C.ulonglong(0), cGoError(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go iih.Run(ctx)

	client := &http.Client{
		Transport: transport,
	}

	httpTransport := &HTTPTransport{iih: iih, client: client, cancel: cancel, refCount: 1}
	httpTransportPool.Set(httpTransportID, httpTransport)

	return C.ulonglong(httpTransportID), nil
}

//export SendHTTPRequest
func SendHTTPRequest(httpTransportID C.ulonglong, method, url, body *C.char) (*C.struct_GoHTTPResponse, *C.struct_GoError) {
	httpTransport, ok := httpTransportPool.Get(uint64(httpTransportID))
	if !ok {
		return cGoInternalServerErrorHTTPResponse(), cGoError(errors.NewPlain("transport not found"))
	}

	httpResp, err := httpTransport.Request(
		C.GoString(method),
		C.GoString(url),
		C.GoString(body))

	if err != nil {
		return cGoInternalServerErrorHTTPResponse(), cGoError(err)
	}

	var headers map[string]string
	if len(httpResp.Headers) > 0 {
		headers = make(map[string]string, len(httpResp.Headers))
		for k, v := range httpResp.Headers {
			headers[k] = strings.Join(v, ",")
		}
	}

	return cGoHTTPResponse(httpResp.StatusCode, httpResp.Version, headers, httpResp.Body), nil
}

//export CloseHTTPTransport
func CloseHTTPTransport(httpTransportID C.ulonglong) {
	if httpTransport, ok := httpTransportPool.GetAndDelete(uint64(httpTransportID)); ok {
		if httpTransport.refCount <= 0 {
			httpTransport.Close()
		}
	}
}

type HTTPTransport struct {
	iih      *istio.IstioIntegrationHandler
	client   *http.Client
	cancel   context.CancelFunc
	refCount uint32
}

type HTTPResponse struct {
	StatusCode int
	Version    string
	Headers    http.Header
	Body       []byte
}

func (t *HTTPTransport) Request(method, url, body string) (*HTTPResponse, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	request, err := http.NewRequestWithContext(context.Background(), method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	response, err := t.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		StatusCode: response.StatusCode,
		Version:    response.Proto,
		Headers:    response.Header,
		Body:       bodyBytes,
	}, nil
}

func (t *HTTPTransport) Close() {
	t.cancel()
}

type httpTransportPoolManager struct {
	pool map[uint64]*HTTPTransport
	mu   sync.RWMutex
}

func (m *httpTransportPoolManager) Get(id uint64) (*HTTPTransport, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	httpTransport, ok := m.pool[id]
	return httpTransport, ok
}

func (m *httpTransportPoolManager) GetAndDelete(id uint64) (*HTTPTransport, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	httpTransport, ok := m.pool[id]
	if ok {
		httpTransport.refCount--
		if httpTransport.refCount <= 0 {
			delete(m.pool, id)
		}
	}
	return httpTransport, ok
}

func (m *httpTransportPoolManager) Set(id uint64, httpTransport *HTTPTransport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pool == nil {
		m.pool = make(map[uint64]*HTTPTransport)
	}
	m.pool[id] = httpTransport
}

func cGoError(err error) *C.struct_GoError {
	goErr := C.create_go_error(C.CString(err.Error()))

	return goErr
}

//export FreeGoError
func FreeGoError(p *C.struct_GoError) {
	ptr := unsafe.Pointer(p)
	if ptr != nil {
		C.free(unsafe.Pointer(p.error_msg))
	}

	C.free(ptr)
}

func cGoHTTPResponse(statusCode int, version string, httpHeaders map[string]string, body []byte) *C.struct_GoHTTPResponse {
	var headers *C.struct_GoHTTPHeaderItem
	count := 0

	if len(httpHeaders) > 0 {
		headers = C.new_http_headers(C.uint(len(httpHeaders)))
		for k, v := range httpHeaders {
			hdr := C.struct_GoHTTPHeaderItem{
				key:   C.CString(k),
				value: C.CString(v),
			}

			C.set_http_header_at(headers, hdr, C.uint(count))
			count++
		}
	}

	httpResponse := C.create_go_http_response(
		C.int(statusCode),
		C.CString(version),
		C.struct_GoHTTPHeaders{
			items: headers,
			size:  C.uint(count),
		},
		C.CString(string(body)),
	)

	return httpResponse
}

func cGoInternalServerErrorHTTPResponse() *C.struct_GoHTTPResponse {
	return cGoHTTPResponse(http.StatusInternalServerError, "", nil, nil)
}

//export FreeGoHTTPResponse
func FreeGoHTTPResponse(p *C.struct_GoHTTPResponse) {
	ptr := unsafe.Pointer(p)
	if ptr != nil {
		C.free(unsafe.Pointer(p.version))
		C.free_http_headers(p.headers.items, p.headers.size)
		C.free(unsafe.Pointer(p.body))
	}

	C.free(ptr)
}

func main() {}
