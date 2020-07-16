FROM ubuntu:xenial

ARG VERSION=7.17.2
ARG ARCH=x86_64
ARG RPCUSERNAME=user
ARG RPCPASSWORD=pass
ARG ROOTDATADIR=/root/

# Set to 1 for running it in testnet mode
ARG TESTNET=0

# Do we want any blockchain pruning to take place? Set to 4096 for a 4GB blockchain prune.
# Alternatively set size=1 to prune with RPC call 'pruneblockchainheight <height>'
# Must have txindex=0 if set
ARG PRUNESIZE=0


RUN apt-get update && apt-get install --yes software-properties-common
RUN apt-get install --yes curl

RUN curl https://github.com/DigiByte-Core/DigiByte/releases/download/v${VERSION}/digibyte-${VERSION}-${ARCH}-linux-gnu.tar.gz -O - | tar xz

RUN mv ./digibyte-${VERSION} /app
RUN chmod +x /app/bin/digibyted
RUN chmod +x /app/bin/digibyte-cli

RUN echo -e "datadir=${ROOTDATADIR}/.digibyte/\n\
server=1\n\
prune=${PRUNESIZE}\n\
maxconnections=300\n\
rpcallowip=127.0.0.1\n\
daemon=1\n\
rpcuser=${RPCUSERNAME}\n\
rpcpassword=$RPCPASSWORD}\n\
txindex=1\n\
addresstype=bech32\n\
# Uncomment below if you need Dandelion disabled for any reason but it is left on by default intentionally\n\
#disabledandelion=1\n\
testnet=${TESTNET}\n\
regtest=1\n" > ${ROOTDATADIR}/.digibyte/digibyte.conf

# Allow Mainnet RPC
EXPOSE 14022

# Allow Testnet RPC
EXPOSE 14023

# Allow Regtest RPC
EXPOSE 18444

COPY run.sh /root/
RUN chmod +x /root/run.sh

# Allow Mainnet RPC
EXPOSE 14022

# Allow Testnet RPC
EXPOSE 14023


ENTRYPOINT ["./root/run.sh"]
