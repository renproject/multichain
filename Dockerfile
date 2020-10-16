FROM ubuntu:bionic

RUN apt update -y
RUN apt install -y mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl wget
RUN apt upgrade -y

RUN wget -c https://golang.org/dl/go1.14.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.14.6.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOSUMDB=off

RUN mkdir -p $(go env GOPATH)
WORKDIR $GOPATH
RUN mkdir -p src/github.com/filecoin-project
WORKDIR $GOPATH/src/github.com/filecoin-project
RUN git clone https://github.com/filecoin-project/filecoin-ffi
WORKDIR $GOPATH/src/github.com/filecoin-project/filecoin-ffi
RUN git checkout a62d00da59d1b0fb
RUN make
RUN go install
