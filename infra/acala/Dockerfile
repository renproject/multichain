FROM ubuntu:bionic

RUN apt-get update && apt-get install --yes --fix-missing software-properties-common curl git clang
RUN apt-get install --yes --fix-missing --no-install-recommends build-essential

# Install Rust
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Clone repository
RUN git clone https://github.com/AcalaNetwork/Acala.git

RUN mv Acala /app
WORKDIR /app

# Make sure submodule.recurse is set to true to make life with submodule easier.
RUN git fetch && git checkout c71dd4e6a5675a0a58a61d7d877a7f02daa3d99f
RUN git config --global submodule.recurse true

# Build
RUN rustup default nightly-2020-09-27
RUN rustup target add wasm32-unknown-unknown --toolchain nightly-2020-09-27
RUN make submodule
RUN make build-full
RUN make build

WORKDIR /
COPY run.sh /root/
RUN chmod +x /root/run.sh

# rpc port
EXPOSE 9933
# ws port
EXPOSE 9944

ENTRYPOINT ["./root/run.sh"]
