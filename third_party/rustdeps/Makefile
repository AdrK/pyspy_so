RUSTC_TARGET ?= `uname -m`-unknown-linux-gnu

.PHONY: all
all:
	wget -nc https://github.com/libunwind/libunwind/releases/download/v1.5/libunwind-1.5.0.tar.gz && \
		tar -zxf libunwind-1.5.0.tar.gz && \
		cd libunwind-1.5.0/ && \
		mkdir -p build/install && \
		cd build && \
		../configure --prefix=${PWD}/libunwind-1.5.0/build/install --disable-minidebuginfo --with-pic --enable-ptrace --disable-tests --disable-documentation && \
		make -j`nproc` && \
		make install && \
		cd -

	cp ./libunwind-1.5.0/build/install/lib/*.a ./

	RUSTFLAGS="-C relocation-model=pic -L${PWD}" cargo build --release --target ${RUSTC_TARGET}
#RUSTFLAGS="-C relocation-model=pic -C target-feature=+crt-static -L${PWD}" cargo build --release --target ${RUSTC_TARGET}

	cp ./target/${RUSTC_TARGET}/release/*.a ./

.PHONY: clean
clean::
	rm -fr ./libunwind*
	rm -fr ./librustdeps*
	rm -fr ./Cargo.lock
	rm -fr ./target/
