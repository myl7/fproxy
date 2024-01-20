SHELL = /bin/sh
.SUFFIXES:
.SUFFIXES: .go
INSTALL = install -DT
INSTALL_PROGRAM = $(INSTALL)
INSTALL_DATA = $(INSTALL) -m 644
prefix = /usr/local
exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin
sysconfdir = $(prefix)/etc

.PHONY: all install install-deb build mkdir

all: build

install: build
	$(INSTALL_PROGRAM) build/fproxy $(DESTDIR)$(bindir)/fproxy
	$(INSTALL_DATA) init/systemd/fproxy.service $(DESTDIR)$(sysconfdir)/systemd/system/fproxy.service
	$(INSTALL_DATA) init/systemd/fproxy $(DESTDIR)$(sysconfdir)/default/fproxy
install-deb: install
	$(INSTALL_PROGRAM) init/debian/postinst $(DESTDIR)/DEBIAN/postinst
	$(INSTALL_DATA) init/debian/control $(DESTDIR)/DEBIAN/control

build: mkdir build/fproxy
build/fproxy: $(wildcard cmd/fproxy/*.go) $(wildcard *.go)
	go build -o $@ $(wildcard cmd/fproxy/*.go)
mkdir:
	@mkdir -p build
