RUSTC_TARGET ?= `uname -m`-unknown-linux-gnu
LIBUNWIND_VERSION ?= "1.5.0"

.PHONY: all
all:
	wget -nc https://github.com/libunwind/libunwind/releases/download/v1.5/libunwind-${LIBUNWIND_VERSION}.tar.gz -O libunwind.tar.gz && \
	tar -zxf libunwind.tar.gz || true

	cd libunwind-${LIBUNWIND_VERSION}/ && \
	mkdir -p install && \
	ls config.status || ./configure --prefix=${PWD}/libunwind-${LIBUNWIND_VERSION}/install --disable-minidebuginfo --with-pic --enable-ptrace --disable-tests --disable-documentation && \
	make -j`nproc` && \
	make install

	cp ./libunwind-${LIBUNWIND_VERSION}/install/lib/*.a ./

#RUSTFLAGS="-C relocation-model=pic -L${PWD}" cargo build --release --target ${RUSTC_TARGET}
	RUSTFLAGS="-C target-feature=+crt-static -L${PWD}" cargo build --release --target ${RUSTC_TARGET}

	cp ./target/${RUSTC_TARGET}/release/*.a ./

.PHONY: clean
clean::
	rm -fr ./libunwind*
	rm -fr ./librustdeps*
	rm -fr ./Cargo.lock
	rm -fr ./target/
