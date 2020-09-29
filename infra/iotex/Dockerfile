FROM golang:1.14-alpine as build
RUN apk add --no-cache make gcc musl-dev linux-headers git curl
ENV GO111MODULE=on

ARG CORE_VERSION="v1.1.1"
RUN git clone --single-branch --branch $CORE_VERSION https://github.com/iotexproject/iotex-core /usr/local/build-core && \
    cd /usr/local/build-core && \
    make clean build && \
    cp ./bin/server /usr/local/bin/iotex-server  && \
    cd / && rm -rf /usr/local/build-core

ARG POLL_META_URI="https://storage.googleapis.com/blockchain-golden/poll.mainnet.tar.gz"
RUN mkdir -p /etc/iotex/meta && \
    curl -L $POLL_META_URI > /etc/iotex/meta/poll.tar.gz

FROM alpine:latest
RUN apk add --no-cache ca-certificates

RUN mkdir -p /etc/iotex/meta
COPY --from=build /usr/local/bin/iotex-server /usr/local/bin
COPY --from=build /etc/iotex/meta/poll.tar.gz /etc/iotex/meta

RUN cd /etc/iotex/meta && \
    tar -xzf poll.tar.gz

VOLUME /data
WORKDIR /data

CMD iotex-server -config-path=/data/etc/config_override.yaml -genesis-path=/data/etc/genesis.yaml -plugin=gateway

EXPOSE 8080/tcp
EXPOSE 14014/tcp
EXPOSE 4689/tcp
