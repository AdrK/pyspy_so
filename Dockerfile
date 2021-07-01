FROM golang:1.15.1-alpine3.12

RUN apk update &&\
    apk add --no-cache build-base openssl-dev musl musl-dev curl make git zstd gcc g++ bash python2 gdb

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN /root/.cargo/bin/rustup target add $(uname -m)-unknown-linux-musl
RUN ln -s /root/.cargo/bin/cargo /usr/bin/cargo

WORKDIR /opt/pyroscope
COPY . .

RUN make clean
RUN make RUSTC_TARGET=x86_64-unknown-linux-musl build-rust-dependencies
RUN make RUSTC_TARGET=x86_64-unknown-linux-musl build-shared
RUN make RUSTC_TARGET=x86_64-unknown-linux-musl build-exe
