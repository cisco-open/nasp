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

package istio

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"github.com/cisco-open/nasp/pkg/ca"
	"github.com/cisco-open/nasp/pkg/istio/discovery"
	"github.com/cisco-open/nasp/pkg/network"
)

type istioHTTPRequestTransport struct {
	transport       http.RoundTripper
	caClient        ca.Client
	discoveryClient discovery.DiscoveryClient
	logger          logr.Logger
}

func NewIstioHTTPRequestTransport(transport http.RoundTripper, caClient ca.Client, discoveryClient discovery.DiscoveryClient, logger logr.Logger) http.RoundTripper {
	return &istioHTTPRequestTransport{
		transport:       transport,
		caClient:        caClient,
		discoveryClient: discoveryClient,
		logger:          logger,
	}
}

func (t *istioHTTPRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var tlsConfig *tls.Config

	if prop, _ := t.discoveryClient.GetHTTPClientPropertiesByHost(context.Background(), req.URL.Host); prop != nil { //nolint:nestif
		t.logger.V(3).Info("discovered overrides", "overrides", prop)
		if endpointAddr, err := prop.Address(); err != nil {
			return nil, err
		} else {
			req.URL.Host = endpointAddr.String()
			t.logger.V(3).Info("address override", "address", req.URL.Host)
		}

		if prop.UseTLS() {
			tlsConfig = t.getTLSconfig()
			tlsConfig.ServerName = prop.ServerName()
			req.URL.Scheme = "https"
			if !strings.Contains(req.URL.Host, ":") {
				req.URL.Host += ":80"
			}
			tlsConfig.NextProtos = []string{
				"istio-http/1.1",
				"istio",
				"http/1.1",
			}
		}
	}

	return network.WrapHTTPTransport(t.transport, network.NewDialerWithTLSConfig(tlsConfig)).RoundTrip(req)
}

func (t *istioHTTPRequestTransport) getTLSconfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(t.caClient.GetCAPem())

	return &tls.Config{
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			cert, err := t.caClient.GetCertificate("", time.Duration(168)*time.Hour)
			if err != nil {
				return nil, err
			}

			return cert.GetTLSCertificate(), nil
		},
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
}
