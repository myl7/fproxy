// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package fproxy

import "net/http"

type Proxy struct {
	c Config
}

func NewProxy(c Config) *Proxy {
	return &Proxy{c}
}

// ServeHTTP implements http.Handler
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

type Config struct {
	HostMap map[string]string
}
