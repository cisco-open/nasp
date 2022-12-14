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

package discovery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"sync"
	"time"

	adsconfig "github.com/cisco-open/nasp/pkg/ads/config"

	"emperror.dev/errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-logr/logr"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/cisco-open/nasp/pkg/ads"
	"github.com/cisco-open/nasp/pkg/ca"
	"github.com/cisco-open/nasp/pkg/environment"
)

type xdsDiscoveryClient struct {
	environment *environment.IstioEnvironment
	caClient    ca.Client
	logger      logr.Logger

	xdsClient ads.Client

	listenerPropertiesContexts   map[string]*listenerPropertiesContext
	tcpClientPropertiesContexts  map[string]*clientPropertiesContext
	httpClientPropertiesContexts map[string]*httpClientPropertiesContext

	mu sync.Mutex
}

type listenerPropertiesContext struct {
	context.Context
	response ads.ListenerPropertiesResponse
}

type clientPropertiesContext struct {
	context.Context
	response ads.ClientPropertiesResponse
}

type httpClientPropertiesContext struct {
	context.Context
	response ads.HTTPClientPropertiesResponse
}

func NewXDSDiscoveryClient(environment *environment.IstioEnvironment, caClient ca.Client, logger logr.Logger) DiscoveryClient {
	return &xdsDiscoveryClient{
		environment: environment,
		caClient:    caClient,
		logger:      logger,

		listenerPropertiesContexts:   map[string]*listenerPropertiesContext{},
		tcpClientPropertiesContexts:  map[string]*clientPropertiesContext{},
		httpClientPropertiesContexts: map[string]*httpClientPropertiesContext{},
	}
}

func (d *xdsDiscoveryClient) Connect(ctx context.Context) error {
	d.logger.Info("get certificate")
	cert, err := d.caClient.GetCertificate(d.environment.GetSpiffeID(), time.Hour*24)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp", d.caClient.GetCAEndpoint())
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{*cert.GetTLSCertificate()},
		RootCAs:            x509.NewCertPool(),
		InsecureSkipVerify: true,
	}
	tlsConfig.RootCAs.AppendCertsFromPEM(d.caClient.GetCAPem())

	md, err := structpb.NewStruct(d.environment.GetNodePropertiesFromEnvironment()["metadata"].(map[string]interface{}))
	if err != nil {
		return err
	}

	clientConfig := &adsconfig.ClientConfig{
		NodeInfo: &adsconfig.NodeInfo{
			// workload id - update if id changed (e.g. workload was restarted)
			Id:          d.environment.GetNodeID(),
			ClusterName: d.environment.GetClusterName(),
			Metadata:    md,
		},
		ManagementServerAddress: addr,
		ClusterID:               d.environment.ClusterID,
		TLSConfig:               tlsConfig,
		SearchDomains:           d.environment.SearchDomains,
	}

	d.logger.Info("connecting to XDS server", "nodeID", clientConfig.NodeInfo.Id, "addr", clientConfig.ManagementServerAddress.String(), "clusterID", d.environment.ClusterID, "network", d.environment.Network)

	d.xdsClient, err = ads.Connect(ctx, clientConfig)
	if err != nil {
		return err
	}

	return nil
}

func (d *xdsDiscoveryClient) GetListenerProperties(ctx context.Context, address string, callbacks ...func(ListenerProperties)) (ListenerProperties, error) {
	if d.xdsClient == nil {
		return nil, errors.New("xds client is not connected")
	}

	err := d.watchForListenerPropertiesResponse(ctx, address, callbacks...)
	if err != nil {
		return nil, err
	}

	var response ads.ListenerPropertiesResponse
	err = d.backoffRetry(ctx, func() error {
		if c, ok := d.listenerPropertiesContexts[address]; ok && c.response != nil {
			response = c.response
			return nil
		}
		return errors.New("could not find listener properties")
	})
	if err != nil {
		return nil, err
	}

	return response.ListenerProperties(), response.Error()
}

func (d *xdsDiscoveryClient) watchForListenerPropertiesResponse(ctx context.Context, address string, callbacks ...func(ListenerProperties)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.listenerPropertiesContexts[address]; ok {
		return nil
	}

	d.listenerPropertiesContexts[address] = &listenerPropertiesContext{
		Context: ctx,
	}

	respChan, err := d.xdsClient.GetListenerProperties(ctx, address)
	if err != nil {
		return err
	}

	go func(ctx context.Context, address string, respChan <-chan ads.ListenerPropertiesResponse) {
		for {
			select {
			case <-ctx.Done():
				return
			case resp, ok := <-respChan:
				if ok {
					d.updateListenerPropertiesResponse(address, resp)
					if resp.Error() == nil {
						for _, f := range callbacks {
							f(resp.ListenerProperties())
						}
					}
				}
			}
		}
	}(ctx, address, respChan)

	return nil
}

func (d *xdsDiscoveryClient) updateListenerPropertiesResponse(address string, response ads.ListenerPropertiesResponse) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if c, ok := d.listenerPropertiesContexts[address]; ok {
		d.logger.Info("listener property update", "address", address, "properties", response.ListenerProperties(), "error", response.Error())
		c.response = response
		return
	}
}

func (d *xdsDiscoveryClient) backoffRetry(ctx context.Context, o backoff.Operation) error {
	backoffProperties := backoff.NewExponentialBackOff()
	backoffProperties.MaxElapsedTime = time.Second * 2

	return backoff.Retry(o, backoff.WithContext(backoffProperties, ctx))
}

func (d *xdsDiscoveryClient) GetTCPClientPropertiesByHost(ctx context.Context, address string, callbacks ...func(ClientProperties)) (ClientProperties, error) {
	if d.xdsClient == nil {
		return nil, errors.New("xds client is not connected")
	}

	err := d.watchForTCPClientProperties(ctx, address, callbacks...)
	if err != nil {
		return nil, err
	}

	var response ads.ClientPropertiesResponse

	err = d.backoffRetry(ctx, func() error {
		if c, ok := d.tcpClientPropertiesContexts[address]; ok && c.response != nil {
			response = c.response
			return nil
		}
		return errors.New("could not find listener properties")
	})
	if err != nil {
		return nil, err
	}

	return response.ClientProperties(), response.Error()
}

func (d *xdsDiscoveryClient) watchForTCPClientProperties(ctx context.Context, address string, callbacks ...func(ClientProperties)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.tcpClientPropertiesContexts[address]; ok {
		return nil
	}

	d.tcpClientPropertiesContexts[address] = &clientPropertiesContext{
		Context: ctx,
	}

	respChan, err := d.xdsClient.GetTCPClientPropertiesByHost(ctx, address)
	if err != nil {
		return err
	}

	go func(ctx context.Context, address string, respChan <-chan ads.ClientPropertiesResponse) {
		for {
			select {
			case <-ctx.Done():
				return
			case resp, ok := <-respChan:
				if ok {
					d.updateTCPClientProperties(address, resp)
					if resp.Error() == nil {
						for _, f := range callbacks {
							f(resp.ClientProperties())
						}
					}
				}
			}
		}
	}(ctx, address, respChan)

	return nil
}

func (d *xdsDiscoveryClient) updateTCPClientProperties(address string, response ads.ClientPropertiesResponse) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if c, ok := d.tcpClientPropertiesContexts[address]; ok {
		d.logger.Info("tcp client property update", "address", address, "properties", response.ClientProperties(), "error", response.Error())
		c.response = response
		return
	}
}

func (d *xdsDiscoveryClient) GetHTTPClientPropertiesByHost(ctx context.Context, address string, callbacks ...func(HTTPClientProperties)) (HTTPClientProperties, error) {
	if d.xdsClient == nil {
		return nil, errors.New("xds client is not connected")
	}

	err := d.watchForHTTPClientProperties(ctx, address, callbacks...)
	if err != nil {
		return nil, err
	}

	var response ads.HTTPClientPropertiesResponse

	err = d.backoffRetry(ctx, func() error {
		if c, ok := d.httpClientPropertiesContexts[address]; ok && c.response != nil {
			response = c.response
			return nil
		}
		return errors.New("could not find http client properties")
	})
	if err != nil {
		return nil, err
	}

	return response.ClientProperties(), response.Error()
}

func (d *xdsDiscoveryClient) watchForHTTPClientProperties(ctx context.Context, address string, callbacks ...func(HTTPClientProperties)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.httpClientPropertiesContexts[address]; ok {
		return nil
	}

	d.httpClientPropertiesContexts[address] = &httpClientPropertiesContext{
		Context: ctx,
	}

	respChan, err := d.xdsClient.GetHTTPClientPropertiesByHost(ctx, address)
	if err != nil {
		return err
	}

	go func(ctx context.Context, address string, respChan <-chan ads.HTTPClientPropertiesResponse) {
		for {
			select {
			case <-ctx.Done():
				return
			case resp, ok := <-respChan:
				if ok {
					d.updateHTTPClientListenerProperties(address, resp)
					if resp.Error() == nil {
						for _, f := range callbacks {
							f(resp.ClientProperties())
						}
					}
				}
			}
		}
	}(ctx, address, respChan)

	return nil
}

func (d *xdsDiscoveryClient) updateHTTPClientListenerProperties(address string, response ads.HTTPClientPropertiesResponse) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if c, ok := d.httpClientPropertiesContexts[address]; ok {
		d.logger.V(3).Info("http client property update", "address", address, "properties", response.ClientProperties(), "error", response.Error())
		c.response = response
		return
	}
}
