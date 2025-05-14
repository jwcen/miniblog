// Copyright 2024 jayvee <jvvcen@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	genericoptions "github.com/onexstack/onexstack/pkg/options"

	"github.com/jwcen/miniblog/internal/pkg/log"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, tlsOptions *genericoptions.TLSOptions,
	handler http.Handler) *HTTPServer {

	var tlsConfig *tls.Config
	if tlsOptions != nil && tlsOptions.UseTLS {
		tlsConfig = tlsOptions.MustTLSConfig()
	}

	return &HTTPServer{
		srv: &http.Server{
			Addr:      httpOptions.Addr,
			TLSConfig: tlsConfig,
			Handler:   handler,
		},
	}
}

func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	// 默认启动 HTTP 服务器
	serveFn := func() error {
		return s.srv.ListenAndServe()
	}

	if s.srv.TLSConfig != nil {
		serveFn = func() error {
			return s.srv.ListenAndServeTLS("", "")
		}
	}

	if err := serveFn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
