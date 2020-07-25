FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common

# Install all of the dependencies that are needed to build nixd and
# nix-cli from source.
RUN apt-get update && apt-get install -y \
    automake \
    bsdmainutils \
    curl \
    g++ \
    git \
    libboost-all-dev \
    libevent-dev \
    libssl-dev \
    libtool \
    libzmq3-dev \
    make \
    openjdk-8-jdk \
    pkg-config \
    zlib1g-dev \
    libminizip-dev \
    cmake \
    libgmp-dev

# Install Berkeley DB 4.8
RUN curl -L http://download.oracle.com/berkeley-db/db-4.8.30.tar.gz | tar -xz -C /tmp && \
    cd /tmp/db-4.8.30/build_unix && \
    ../dist/configure --enable-cxx --includedir=/usr/include/bdb4.8 --libdir=/usr/lib && \
    make && make install && \
    cd / && rm -rf /tmp/db-4.8.30

# Clone the repository and build.
RUN git clone https://github.com/nixplatform/nixcore && \
    cd nixcore && \
    ./autogen.sh && \
    ./configure --without-gui --without-miniupnpc && \
    make && \
    make install

COPY nix.conf /root/.nix/nix.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

# Regtest network port
EXPOSE 16216

ENTRYPOINT ["./root/run.sh"]
