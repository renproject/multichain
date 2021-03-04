FROM bitcoinsv/bitcoin-sv:latest
ADD file:9b8be2b52ee0fa31da1b6256099030b73546253a57e94cccb24605cd888bb74d in / 
CMD ["bash"]
/bin/sh -c groupadd -r bitcoin && useradd -r -m -g bitcoin bitcoin
/bin/sh -c set -ex 	&& apt-get update 	&& apt-get install -qq --no-install-recommends ca-certificates dirmngr gosu gpg wget libatomic1 	&& rm -rf /var/lib/apt/lists/*
ENV BITCOIN_VERSION=1.0.7
ENV BITCOIN_URL=http://download.bitcoinsv.io/bitcoinsv/1.0.7/bitcoin-sv-1.0.7-x86_64-linux-gnu.tar.gz
ENV BITCOIN_SHA256=6785f117b282aee4f81fb06d1fb5fb900306ca5c5a91ad85619f68da20f68445
/bin/sh -c set -ex 	&& cd /tmp 	&& wget -qO bitcoin.tar.gz "$BITCOIN_URL" 	&& echo "$BITCOIN_SHA256 bitcoin.tar.gz" | sha256sum -c - 	&& tar -xzvf bitcoin.tar.gz -C /usr/local --strip-components=1 --exclude=*-qt 	&& rm -rf /tmp/*
ENV BITCOIN_DATA=/data
/bin/sh -c mkdir "$BITCOIN_DATA" 	&& chown -R bitcoin:bitcoin "$BITCOIN_DATA" 	&& ln -sfn "$BITCOIN_DATA" /home/bitcoin/.bitcoin 	&& chown -h bitcoin:bitcoin /home/bitcoin/.bitcoin
VOLUME [/data]
COPY file:8032f05a589c63c27eb4f80ad9517940c77a40820daa5314e1d0d35467654360 in /entrypoint.sh 
ENTRYPOINT ["/entrypoint.sh"]
EXPOSE 18332 18333 8332 8333 9332 9333
CMD ["bitcoind"]
