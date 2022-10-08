// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net/http"

	"github.com/myl7/fproxy"
)

func main() {
	scheme := flag.String("scheme", "https", "scheme")
	host := flag.String("host", "", "host")
	pPrefix := flag.String("pPrefix", "", "path prefix")
	cacheDir := flag.String("cacheDir", "", "cache dir")
	listen := flag.String("listen", ":8000", "listen")
	flag.Parse()

	if *host == "" {
		panic("host is invalid")
	}

	p := fproxy.NewProxy(fproxy.Config{
		URLTransform:   fproxy.URLTransformPrefix(*scheme, *host, *pPrefix, false),
		URLToLocalPath: fproxy.URLToLocalPathPrefix(*cacheDir),
	})
	http.ListenAndServe(*listen, p)
}
