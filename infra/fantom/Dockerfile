FROM ubuntu:groovy

# Install dependencies
RUN apt-get update -y
RUN apt-get install -y build-essential git wget

# Install Go
RUN wget -c https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Build Opera
WORKDIR /app
RUN git clone https://github.com/Fantom-foundation/go-opera.git
WORKDIR /app/go-opera
RUN git checkout release/1.0.0-rc.4
RUN make
ENV PATH=$PATH:/app/go-opera/build

# Expose the default port of the JSON-RPC server
EXPOSE 18545

ENTRYPOINT [ "opera", "--fakenet", "1/1", "--http", "--http.api=eth,net", "--http.vhosts", "*", "--http.addr", "0.0.0.0" ]
