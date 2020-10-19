FROM ubuntu:xenial

RUN apt-get update && apt-get install --yes --fix-missing software-properties-common curl git clang
RUN apt-get install --yes --fix-missing --no-install-recommends build-essential

# Install Rust
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
RUN rustup default nightly
RUN rustup update

# Clone repository
RUN git clone https://github.com/AcalaNetwork/Acala.git

RUN mv Acala /app
WORKDIR /app

# Make sure submodule.recurse is set to true to make life with submodule easier.
RUN git fetch
RUN git pull
RUN git config --global submodule.recurse true

# Build
RUN make init
RUN make build

WORKDIR /
COPY run.sh /root/
RUN chmod +x /root/run.sh

# rpc port
EXPOSE 9933
# ws port
EXPOSE 9944

ENTRYPOINT ["./root/run.sh"]
