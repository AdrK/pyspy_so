#                 _
#                | |
#  _ __ _   _ ___| |_
# | '__| | | / __| __|
# | |  | |_| \__ \ |_
# |_|   \__,_|___/\__|

FROM alpine:3.12 as rust-builder

RUN apk update &&\
    apk add --no-cache git gcc g++ make build-base openssl-dev musl musl-dev curl

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN /root/.cargo/bin/rustup target add $(uname -m)-unknown-linux-musl

RUN wget -nc https://github.com/libunwind/libunwind/releases/download/v1.5/libunwind-1.5.0.tar.gz && \
		tar -zxf libunwind-1.5.0.tar.gz && \
		cd libunwind-1.5.0/ && \
		mkdir -p build/install && \
		cd build && \
		../configure --disable-minidebuginfo --with-pic --enable-ptrace --disable-tests --disable-documentation && \
		make install -j`nproc`

COPY third_party/rustdeps /opt/rustdeps

WORKDIR /opt/rustdeps

RUN RUSTFLAGS="-C relocation-model=pic -C target-feature=+crt-static" /root/.cargo/bin/cargo build --release --target `uname -m`-unknown-linux-musl
RUN mv /opt/rustdeps/target/$(uname -m)-unknown-linux-musl/release/librustdeps.a /opt/rustdeps/librustdeps.a


#              _
#             | |
#   __ _  ___ | | __ _ _ __   __ _
#  / _` |/ _ \| |/ _` | '_ \ / _` |
# | (_| | (_) | | (_| | | | | (_| |
#  \__, |\___/|_|\__,_|_| |_|\__, |
#   __/ |                     __/ |
#  |___/                     |___/

FROM golang:1.15.1-alpine3.12 as go-builder

RUN apk add --no-cache make git zstd gcc g++ libc-dev musl-dev bash python2 gdb

WORKDIR /opt/pyroscope

RUN mkdir -p /opt/pyroscope/third_party/rustdeps/target/release
COPY --from=rust-builder /opt/rustdeps/librustdeps.a /opt/pyroscope/third_party/rustdeps/librustdeps.a
COPY third_party/rustdeps/pyspy.h /opt/pyroscope/third_party/rustdeps/pyspy.h
COPY Makefile ./
COPY go.mod go.sum main.go ./
COPY pkg ./pkg

RUN make build-shared
RUN make build-exe
