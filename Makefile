ALL_SPIES ?= "pyspy"
ENABLED_SPIES ?= "pyspy"
RUSTC_TARGET ?= `uname -m`-unknown-linux-gnu
GOBUILD=go build -trimpath -a -tags $(ENABLED_SPIES)
#GOBUILD=go build  -a -tags $(ENABLED_SPIES) -trimpath -gcflags '-N -l'

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
		make -j`nproc` && \
		make install && \
		cd -

	cp ./third_party/rustdeps/libunwind-1.5.0/build/install/lib/*.a ./third_party/rustdeps/

	cd third_party/rustdeps && \
	RUSTFLAGS="-C relocation-model=pic -L${PWD}/third_party/rustdeps/" cargo build --release --target ${RUSTC_TARGET}
#RUSTFLAGS="-C relocation-model=pic -C target-feature=+crt-static -L${PWD}/third_party/rustdeps/" cargo build --release --target ${RUSTC_TARGET}

	cp ./third_party/rustdeps/target/${RUSTC_TARGET}/release/*.a ./third_party/rustdeps/

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
clean:
	rm -fr third_party/rustdeps/libunwind*
	rm -fr third_party/rustdeps/target/
	rm -fr main libpyspy.a libpyspy.so libpyspy.h