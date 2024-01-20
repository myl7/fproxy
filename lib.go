// Copyright (C) myl7
// SPDX-License-Identifier: Apache-2.0

package fproxy

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type Proxy struct {
	urlPathPrefix string
	cacheDir      string
}

func NewProxy(urlPathPrefix, cacheDir string) *Proxy {
	if urlPathPrefix[0] != '/' {
		urlPathPrefix = "/" + urlPathPrefix
	}
	if urlPathPrefix[len(urlPathPrefix)-1] != '/' {
		urlPathPrefix += "/"
	}

	return &Proxy{
		urlPathPrefix: urlPathPrefix,
		cacheDir:      cacheDir,
	}
}

func (p Proxy) cachePath(urlPath string) string {
	if !strings.HasPrefix(urlPath, p.urlPathPrefix) {
		return ""
	}

	return path.Join(p.cacheDir, strings.TrimPrefix(urlPath, p.urlPathPrefix))
}

func (p Proxy) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPath := r.URL.Path
	if urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}

	cachePath := p.cachePath(urlPath)
	if cachePath == "" {
		p.handleForward(w, r)
		return
	}

	if p.cached(cachePath) {
		// TODO: cachePath security
		http.ServeFile(w, r, cachePath)
	} else {
		p.handleCachedForward(w, r, cachePath)
	}
}

func (p Proxy) handleForward(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (p Proxy) cached(cachePath string) bool {
	if _, err := os.Stat(cachePath); err == nil {
		return true
	}
	return false
}

func (p Proxy) handleCachedForward(w http.ResponseWriter, r *http.Request, cachePath string) {
	d := path.Dir(cachePath)
	if err := os.MkdirAll(d, 0777); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Create(cachePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	resp, err := http.Get(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	}

	w.WriteHeader(http.StatusOK)
	defer resp.Body.Close()
	teeR := io.TeeReader(resp.Body, f)
	io.Copy(w, teeR)
}
