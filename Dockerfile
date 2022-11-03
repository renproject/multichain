FROM golang

# doing all updates and installs in a single step and removing the apt cache helps reduce the image size
RUN apt-get update && \
    apt-get install -y \
    mesa-opencl-icd \
    ocl-icd-opencl-dev \
    libssl-dev \
    libudev-dev \
    hwloc \
    libhwloc-dev \
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
RUN git checkout 7912389334e347bbb2eac0520c836830875c39de
RUN make
RUN go install

COPY ../solana-ffi/* $GOPATH/src/github.com/renproject/