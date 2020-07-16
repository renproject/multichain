FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common

# Install all of the dependencies that are needed to build digibyted and
# digibyte-cli from source.
RUN apt-get install --yes build-essential libtool libboost-all-dev autotools-dev automake pkg-config libssl-dev libevent-dev bsdmainutils wget tar

# Install BerkleyDB 4.8
RUN wget https://download.oracle.com/berkeley-db/db-4.8.30.NC.tar.gz
RUN tar -zxvf db-4.8.30.NC.tar.gz && \
    cd db-4.8.30.NC/build_unix && \
    ../dist/configure --enable-cxx --disable-shared --with-pic && \
    make install && \
    cd ../../

# Install git so that we can clone the digibyte repository.
RUN apt-get install --yes git

# Clone the repository and build.
RUN git clone https://github.com/digibyte-core/digibyte && \
    cd digibyte && \
    ./autogen.sh && \
    export CPATH=/usr/local/BerkeleyDB.4.8/include && \
    export LIBRARY_PATH=/usr/local/BerkeleyDB.4.8/lib && \
    export LD_LIBRARY_PATH=/usr/local/BerkeleyDB.4.8/lib && \
    ./configure --without-gui --without-miniupnpc && \
    make && \
    make install

COPY digibyte.conf /root/.digibyte/digibyte.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 18132

ENTRYPOINT ["./root/run.sh"]
