// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"github.com/myl7/fproxy"
)

func main() {
	p := fproxy.NewProxy(fproxy.Config{})
	http.ListenAndServe(":8000", p)
}
