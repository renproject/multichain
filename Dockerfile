FROM ubuntu:bionic

RUN apt update -y
RUN apt install -y mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl wget
RUN apt upgrade -y

RUN wget -c https://golang.org/dl/go1.14.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.14.6.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

RUN export GOROOT=/usr/local/go && \
    export GOPATH=$HOME/go && \
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN mkdir -p $(go env GOPATH) && \
	cd $(go env GOPATH) && \
	ls && \
	mkdir -p src/github.com/filecoin-project && \
	cd src/github.com/filecoin-project && \
	git clone https://github.com/filecoin-project/filecoin-ffi && \
	cd filecoin-ffi && \
	git checkout a62d00da59d1b0fb && \
	make && \
	go install

ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV GOPRIVATE=github.com/renproject/darknode
