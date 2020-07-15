FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common

# Install all of the dependencies that are needed to build dogecoind and
# dogecoin-cli from source.
RUN apt-get install --yes build-essential libtool libboost-all-dev autotools-dev automake pkg-config libssl-dev libevent-dev bsdmainutils wget tar

# Install BerkleyDB 5.1
RUN wget http://download.oracle.com/berkeley-db/db-5.1.29.tar.gz
RUN tar -zxvf db-5.1.29.tar.gz && \
    cd db-5.1.29/build_unix && \
    ../dist/configure --enable-cxx --disable-shared --with-pic && \
    make install && \
    cd ../../

# Install git so that we can clone the dogecoin repository.
RUN apt-get install --yes git

# Clone the repository and build.
RUN git clone https://github.com/dogecoin/dogecoin && \
    cd dogecoin && \
    ./autogen.sh && \
    export CPATH=/usr/local/BerkeleyDB.5.1/include && \
    export LIBRARY_PATH=/usr/local/BerkeleyDB.5.1/lib && \
    export LD_LIBRARY_PATH=/usr/local/BerkeleyDB.5.1/lib && \
    ./configure --without-gui --without-miniupnpc && \
    make && \
    make install

COPY dogecoin.conf /root/.dogecoin/dogecoin.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 18332

ENTRYPOINT ["./root/run.sh"]