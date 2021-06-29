GOBUILD=go build -trimpath

ALL_SPIES ?= "pyspy"
ENABLED_SPIES ?= "pyspy"

ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH || true)
	export GOPATH
endif

.PHONY: all
all: build

.PHONY: build-rust-dependencies
build-rust-dependencies:
	cd third_party/rustdeps && \
	wget -nc https://github.com/libunwind/libunwind/releases/download/v1.5/libunwind-1.5.0.tar.gz && \
		tar -zxf libunwind-1.5.0.tar.gz && \
		cd libunwind-1.5.0/ && \
		mkdir -p build/install && \
		cd build && \
		../configure --prefix=${PWD}/third_party/rustdeps/libunwind-1.5.0/build/install --disable-minidebuginfo --enable-ptrace --disable-tests --disable-documentation && \
		make install && \
		cd -

	cd third_party/rustdeps && \
	RUSTFLAGS="-C target-feature=+crt-static -L${PWD}/third_party/rustdeps/" cargo build --release --target `uname -m`-unknown-linux-musl

.PHONY: build
build:
	$(GOBUILD) -a -tags $(ENABLED_SPIES) -ldflags "-s" -buildmode=c-shared -o libpyspy.so
