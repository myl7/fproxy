#!/usr/bin/env bash
set -euo pipefail

make DESTDIR=build/deb/fproxy prefix=/ exec_prefix=/usr install-deb
dpkg-deb --build build/deb/fproxy
