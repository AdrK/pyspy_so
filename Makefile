ALL_SPIES ?= "pyspy"
ENABLED_SPIES ?= "pyspy"
RUSTC_TARGET ?= `uname -m`-unknown-linux-gnu
GOBUILD=go build -tags $(ENABLED_SPIES) -trimpath
#GOBUILD=go build -tags $(ENABLED_SPIES) -trimpath -gcflags '-N -l'

ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH || true)
	export GOPATH
endif

include ./third_party/rustdeps/Makefile

.PHONY: build-rust-dependencies
build-rust-dependencies:
	cd ./third_party/rustdeps/ && ${MAKE}

.PHONY: build-shared
build-shared:
	$(GOBUILD) -buildmode=c-shared -o libpyspy.so

.PHONY: build-static
build-static:
	$(GOBUILD) -buildmode=c-archive -o libpyspy.a

.PHONY: build-exe
build-exe:
	$(GOBUILD) -o main

.PHONY: clean
clean::
	rm -fr main libpyspy.a libpyspy.so libpyspy.h
	rm -fr pyspy_pyapi/libpyspy.so pyspy_pyapi/libpyspy.h
	cd ./third_party/rustdeps/ && ${MAKE} clean

.PHONY: all
all: build-rust-dependencies build-shared build-static build-exe clean
