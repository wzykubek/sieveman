BIN := sieveman
MODULE := $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")

BUILD_FLAGS = -a -ldflags "-s -w -X $(MODULE)/cmd.version=$(VERSION)" -trimpath

DESTDIR =
PREFIX = /usr/local

.PHONY: all clean fmt vet test build completions docs install install-bin install-completions install-docs uninstall

all: clean build completions

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
	go build $(BUILD_FLAGS) -o dist/$(BIN) $(MODULE)

completions:
	mkdir -p dist/completions
	for sh in bash zsh fish powershell; do \
		dist/$(BIN) completion $$sh > dist/completions/$$sh; \
	done

docs:
	mkdir -p docs/man
	go run ./tools/build-man > ./docs/man/sieveman.1

install: install-bin install-completions install-docs
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

install-docs: docs
	install -d $(DESTDIR)$(PREFIX)/share/man/man1
	gzip -9 < docs/man/$(BIN).1 > $(DESTDIR)$(PREFIX)/share/man/man1/$(BIN).1.gz

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/licenses/$(BIN)/LICENSE
	rm -f $(DESTDIR)$(PREFIX)/share/bash-completion/completions/$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)
	rm -f $(DESTDIR)$(PREFIX)/share/fish/vendor_completions.d/$(BIN).fish
	rm -f $(DESTDIR)$(PREFIX)/share/man/man1/$(BIN).1.gz
