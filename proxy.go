// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package fproxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Proxy struct {
	c Config
}

func NewProxy(c Config) *Proxy {
	return &Proxy{c}
}

// ServeHTTP implements http.Handler
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	u, ok := p.c.URLTransform(*req.URL)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if req.Method != http.MethodGet {
		http.Redirect(w, req, u.String(), http.StatusTemporaryRedirect)
		return
	}

	_, ok = req.Header["Range"]
	if ok {
		p.forwardRangeGet(w, req, u)
		return
	}

	p.forwardGet(w, req, u)
}

func (p *Proxy) forwardGet(w http.ResponseWriter, req *http.Request, u url.URL) {
	resp, err := http.Get(u.String())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	fp := p.c.URLToLocalPath(u)
	fpDir := path.Dir(fp)
	err = os.MkdirAll(fpDir, 0777)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	f, err := os.Create(fp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer f.Close()

	rd := io.TeeReader(resp.Body, f)
	io.Copy(w, rd)
}

func (p *Proxy) forwardRangeGet(w http.ResponseWriter, req *http.Request, u url.URL) {
	http.Error(w, "Unimplemented", http.StatusInternalServerError)
}

type Config struct {
	URLTransform   func(u url.URL) (url.URL, bool)
	URLToLocalPath func(u url.URL) string
}

// pPrefixInsert If true, join pPrefix with path, which inserts it as prefix. If false, trim it from path prefix.
func URLTransformPrefix(scheme, newH, pPrefix string, pPrefixInsert bool) func(u url.URL) (url.URL, bool) {
	return func(u url.URL) (url.URL, bool) {
		u.Scheme = scheme
		u.Host = newH

		if pPrefixInsert {
			u.Path = path.Clean(path.Join("/", pPrefix, u.Path))
		} else if pPrefix != "" && pPrefix != "/" {
			p := path.Clean(path.Join("/", u.Path))
			pp := strings.TrimSuffix(path.Clean(path.Join("/", pPrefix)), "/")
			if strings.HasPrefix(p, pp) && (len(pp) == len(p) || p[len(pp)] == '/') {
				u.Path = p[len(pp):]
			}
		}

		return u, true
	}
}

func URLToLocalPathPrefix(pPrefix string) func(u url.URL) string {
	return func(u url.URL) string {
		pp := pPrefix
		if pp == "" {
			pp = "."
		}
		return path.Clean(path.Join(pp, u.Path))
	}
}
