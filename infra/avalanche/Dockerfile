FROM ubuntu:groovy

RUN apt update -y
RUN apt-get install -y wget curl build-essential

# INSTALL GO AND RUST
RUN wget -c https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Fetch avalanchego
WORKDIR /app
RUN wget https://github.com/ava-labs/avalanchego/releases/download/v1.4.5/avalanchego-linux-amd64-v1.4.5.tar.gz
RUN tar -xzf avalanchego-linux-amd64-v1.4.5.tar.gz
RUN rm avalanchego-linux-amd64-v1.4.5.tar.gz
ENV PATH=$PATH:/app/avalanchego-v1.4.5

# COPY run script
COPY run.sh run.sh
RUN chmod +x run.sh

# RPC SERVER PORT
EXPOSE 9650

# RUN
ENTRYPOINT [ "./run.sh" ]
