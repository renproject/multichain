FROM golang

# doing all updates and installs in a single step and removing the apt cache helps reduce the image size
RUN apt-get update && \
        apt-get install -y \
        mesa-opencl-icd \
        ocl-icd-opencl-dev \
        libssl-dev \
        libudev-dev \
        gcc \
        git \
        bzr \
        jq \
        pkg-config \
        curl \
        wget && \
        apt-get upgrade -y && \
        rm -rf /var/lib/apt/lists/*

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
RUN git checkout a62d00da59d1b0fb35f3a4ae854efa9441af892d
RUN make
RUN go install

WORKDIR $GOPATH
RUN go get -u github.com/xlab/c-for-go
RUN mkdir -p src/github.com/renproject
WORKDIR $GOPATH/src/github.com/renproject
RUN git clone https://github.com/renproject/solana-ffi
WORKDIR $GOPATH/src/github.com/renproject/solana-ffi
RUN git checkout df7838d724f5eaf262ed77ed93b35b3f1f652bd3
RUN make clean && make
RUN go install ./...
