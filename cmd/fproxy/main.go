// Copyright (C) myl7
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/myl7/fproxy"
)

func main() {
	addr := flag.String("addr", ":8000", "Address to listen on, e.g., :8000 or 0.0.0.0:8000")
	urlPathPrefixF := flag.String("urlPathPrefix", "", "URL path prefix whose all subpaths are cached")
	cacheDirF := flag.String("cacheDir", "", "Directory to store cached files")
	flag.Parse()
	if *urlPathPrefixF == "" || *cacheDirF == "" {
		flag.Usage()
		return
	}

	urlPathPrefix := path.Clean(*urlPathPrefixF)
	cacheDir := path.Clean(*cacheDirF)
	if strings.Contains(urlPathPrefix, "..") || strings.Contains(cacheDir, "..") {
		fmt.Fprintln(os.Stderr, "urlPathPrefix and cacheDir must not contain '..'")
		os.Exit(1)
	}

	proxy := fproxy.NewProxy(urlPathPrefix, cacheDir)
	http.HandleFunc("/", proxy.Handle)
	http.ListenAndServe(*addr, nil)
}
