// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package fproxy

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
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

	var isRangeReq bool
	var rangeStart int64
	rangeFields, isRangeReq := req.Header["Range"]
	if isRangeReq {
		m := regexp.MustCompile(`^bytes=(\d+)-\d*$`).FindStringSubmatch(rangeFields[0])
		if m == nil {
			http.Error(w, "", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		var err error
		rangeStart, err = strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			panic(err)
		}

		// http.Error(w, "Unimplemented", http.StatusInternalServerError)
		// return
	}

	subReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		panic(err)
	}

	if isRangeReq {
		subReq.Header.Add("Range", rangeFields[0])
	}

	resp, err := http.DefaultClient.Do(subReq)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	isRangeResp := resp.StatusCode == http.StatusPartialContent

	if !isRangeResp && resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			panic(err)
		}
		return
	}

	fp := p.c.URLToLocalPath(u)
	fpDir := path.Dir(fp)
	err = os.MkdirAll(fpDir, 0777)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Seek(rangeStart, io.SeekStart)

	bufF := bufio.NewWriter(f)
	defer bufF.Flush()

	w.Header().Add("Access-Control-Allow-Origin", "*")

	// TODO: Handle header only

	rd := io.TeeReader(resp.Body, bufF)
	n, err := io.Copy(w, rd)
	if err != nil {
		panic(err)
	}

	if isRangeResp {
		// TODO: Save n
		n += 1
	}
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
