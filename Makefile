PKG := $(shell cat go.mod | awk 'NR==1{print $$2}')
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS :="
LDFLAGS += -X ${PKG}/version.version=$(VERSION)
LDFLAGS += -X ${PKG}/version.buildDate=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS += -X ${PKG}/version.gitCommit=$(COMMIT)
LDFLAGS += -s -w
LDFLAGS +="

.PHONY: all
all: fmtgo docker-build

.PHONY: fmtgo
fmtgo:
	@echo "+ $@"
	@./scripts/gofmt.sh

docker-build:
	@echo "+ $@"
	docker build --build-arg LDFLAGS=${LDFLAGS} -t bennu/identifier-producer .
