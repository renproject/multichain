FROM ubuntu:groovy

# Install dependencies
RUN apt-get update -y
RUN apt-get install -y build-essential git wget

# Install Go
RUN wget -c https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Build lachesis
WORKDIR /app
RUN git clone https://github.com/Fantom-foundation/go-lachesis.git
WORKDIR /app/go-lachesis
RUN git checkout tags/v0.7.0-rc.1 -b lachesis-v7rc1
RUN make build
ENV PATH=$PATH:/app/go-lachesis/build

# Expose the default port of the JSON-RPC server
EXPOSE 18545

ENTRYPOINT [ "lachesis", "--fakenet", "1/1", "--rpc", "--rpcvhosts", "*", "--rpcaddr", "0.0.0.0" ]
