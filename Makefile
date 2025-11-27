BIN := sieveman
MODULE := $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")

LDFLAGS = -linkmode=external -X $(MODULE)/cmd.version=$(VERSION)

DESTDIR =
PREFIX = /usr/local

.PHONY: all clean fmt vet test build completions install install-bin install-completions uninstall

all: clean build

clean:
	rm -rf dist

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p dist
	go build \
		-trimpath \
		-mod=readonly \
		-modcacherw \
		-buildmode=pie \
		-ldflags "$(LDFLAGS)" \
		-o dist/$(BIN) \
		./main.go

completions: build
	mkdir -p dist/completions
	for sh in bash zsh fish powershell; do \
		dist/$(BIN) completion $$sh > dist/completions/$$sh; \
	done

install: install-bin
	install -Dm644 LICENSE $(DESTDIR)$(PREFIX)/share/licenses/$(BIN)/LICENSE

install-bin:
	install -Dm755 dist/$(BIN) $(DESTDIR)$(PREFIX)/bin/$(BIN)

install-completions:
	install -Dm644 dist/completions/bash \
		$(DESTDIR)$(PREFIX)/share/bash-completion/completions/$(BIN)
	install -Dm644 dist/completions/zsh \
		$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)
	install -Dm644 dist/completions/fish \
		$(DESTDIR)$(PREFIX)/share/fish/vendor_completions.d/$(BIN).fish

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/licenses/$(BIN)/LICENSE
	rm -f $(DESTDIR)$(PREFIX)/share/bash-completion/completions/$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/fish/vendor_completions.d/$(BIN).fish
