FROM golang:latest

RUN apt-get -y update && \
    apt-get -y install curl make git zstd gcc g++ bash python2 gdb vim

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN ln -s /root/.cargo/bin/cargo /usr/bin/cargo && ln -s /root/.cargo/bin/rustup /usr/bin/rustup

WORKDIR /opt/pyroscope
COPY . .

RUN make clean
RUN make RUSTC_TARGET=x86_64-unknown-linux-gnu build-rust-dependencies
RUN make RUSTC_TARGET=x86_64-unknown-linux-gnu build-shared
RUN make RUSTC_TARGET=x86_64-unknown-linux-gnu build-exe
