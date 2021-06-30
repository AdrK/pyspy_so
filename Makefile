GOBUILD=go build -trimpath -gcflags '-N -l'

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
		../configure --prefix=${PWD}/third_party/rustdeps/libunwind-1.5.0/build/install --disable-minidebuginfo --with-pic --enable-ptrace --disable-tests --disable-documentation && \
		make install -j`nproc` && \
		cd -

	cp ./third_party/rustdeps/libunwind-1.5.0/build/install/lib/*.a ./third_party/rustdeps/

	cd third_party/rustdeps && \
	RUSTFLAGS="-C relocation-model=pic -C target-feature=+crt-static -L${PWD}/third_party/rustdeps/" cargo build --release --target `uname -m`-unknown-linux-musl

	cp ./third_party/rustdeps/target/`uname -m`-unknown-linux-musl/release/*.a ./third_party/rustdeps/

.PHONY: build-shared
build-shared:
	$(GOBUILD) -a -tags $(ENABLED_SPIES) -buildmode=c-shared -o libpyspy.so

.PHONY: build-exe
build-exe:
	$(GOBUILD) -a -tags $(ENABLED_SPIES) -o main
