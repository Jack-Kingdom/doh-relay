.EXPORT_ALL_VARIABLES:
.PHONY: build bdist dep

GO111MODULE=on

ENTRY = main.go
OUTPUT = bin/doh-relay
VERSION = `git describe --tags`
BUILD = `date +%FT%T%z`
COMPILER = `go version`
LDFLAGS = -ldflags "-s -w -X main.version=${VERSION} -X main.build=${BUILD}"
RABBIT_DEBUG = true

build:
	CGO_ENABLED=0 && go build ${LDFLAGS} -o ${OUTPUT} ${ENTRY}

bdist:
	if [ -d "bdist" ]; then trash bdist; fi
	mkdir -p bdist
	export GOOS=linux GOARCH=amd64 && make build && make pack

dep:
	go mod tidy
