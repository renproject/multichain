FROM golang

RUN apt update -y
RUN apt install -y mesa-opencl-icd ocl-icd-opencl-dev libssl-dev libudev-dev gcc git bzr jq pkg-config curl wget
RUN apt upgrade -y

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

RUN mkdir -p $(go env GOPATH)
WORKDIR $GOPATH
RUN mkdir -p src/github.com/filecoin-project
WORKDIR $GOPATH/src/github.com/filecoin-project
RUN git clone https://github.com/filecoin-project/filecoin-ffi
WORKDIR $GOPATH/src/github.com/filecoin-project/filecoin-ffi
RUN git checkout a62d00da59d1b0fb
RUN make
RUN go install

WORKDIR $GOPATH
RUN go get -u github.com/xlab/c-for-go
RUN mkdir -p src/github.com/renproject
WORKDIR $GOPATH/src/github.com/renproject
RUN git clone https://github.com/renproject/solana-ffi
WORKDIR $GOPATH/src/github.com/renproject/solana-ffi
RUN git checkout 44840392296fa690cd777a55dce19fd4844c1559
RUN make clean && make
RUN go install ./...
