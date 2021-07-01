FROM golang:latest AS rust-builder

RUN apt-get -y update && \
    apt-get -y install curl make gcc

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN ln -s /root/.cargo/bin/cargo /usr/bin/cargo && ln -s /root/.cargo/bin/rustup /usr/bin/rustup

WORKDIR /opt/pyroscope/rustdeps
COPY ./third_party/rustdeps/ ./
RUN make clean
RUN make RUSTC_TARGET=x86_64-unknown-linux-gnu all
RUN cp ./target/x86_64-unknown-linux-gnu/release/*.a /opt/pyroscope/rustdeps/

FROM golang:latest AS go-builder

WORKDIR /opt/pyroscope
COPY --from=rust-builder /opt/pyroscope/rustdeps/*.a ./third_party/rustdeps/
COPY --from=rust-builder /opt/pyroscope/rustdeps/*.h ./third_party/rustdeps/
COPY . .
RUN make RUSTC_TARGET=x86_64-unknown-linux-gnu build-shared build-static build-exe
