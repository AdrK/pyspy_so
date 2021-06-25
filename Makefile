GOBUILD=go build -trimpath

ifeq ("$(shell go env GOARCH || true)", "arm64")
	# this makes it work better on M1 machines
	GODEBUG=asyncpreemptoff=1
endif

ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH || true)
	export GOPATH
endif

.PHONY: all
all: build

.PHONY: build
build:
	$(GOBUILD) -buildmode=c-shared -o libpyspy.so
