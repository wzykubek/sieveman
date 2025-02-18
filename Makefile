BIN=sieveman

VERSION=$(shell git describe --tags --always --dirty | sed 's/^v//')
LDFLAGS=-linkmode=external -X go.wzykubek.xyz/${BIN}/cmd.version=${VERSION}

DESTDIR=
PREFIX=/usr/local

.PHONY: all
all: | clean build completions

.PHONY: clean
clean:
	rm -rf dist/*

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build: $(shell find . -name "*.go" -type f)
	mkdir -p dist
	go build \
		-trimpath \
		-mod=readonly \
		-modcacherw \
		-buildmode=pie \
		-ldflags '${LDFLAGS}' \
		-o dist/${BIN} \
		./main.go

.PHONY: completions
completions:
	mkdir -p dist/completions
	for shell in bash zsh fish powershell; do \
		dist/${BIN} completion $$shell > dist/completions/$$shell; \
	done

.PHONY: install
install:
	install -Dm755 dist/${BIN} ${DESTDIR}${PREFIX}/bin/${BIN}
	install -Dm644 dist/completions/bash \
		${DESTDIR}${PREFIX}/share/bash-completion/completions/${BIN}
	install -Dm644 dist/completions/zsh \
		${DESTDIR}${PREFIX}/share/zsh/site-functions/_${BIN}
	install -Dm644 dist/completions/fish \
		${DESTDIR}${PREFIX}/share/fish/vendor_completions.d/${BIN}.fish
	install -Dm644 LICENSE ${DESTDIR}${PREFIX}/share/licenses/${BIN}/LICENSE

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/${BIN}
	rm -f ${DESTDIR}${PREFIX}/share/bash-completion/completions/${BIN}
	rm -f ${DESTDIR}${PREFIX}/share/zsh/site-functions/_${BIN}
	rm -f ${DESTDIR}${PREFIX}/share/fish/vendor_completions.d/${BIN}.fish
