FROM ubuntu:xenial

# Install bitcoind-abc.
RUN apt-get update && apt-get install --yes software-properties-common && \
add-apt-repository --yes ppa:bitcoin-abc/ppa && apt-get update && \
apt-get install --yes bitcoind

COPY bitcoin.conf /root/.bitcoin/
COPY run.sh /root/
RUN chmod +x /root/run.sh

EXPOSE 18443

ENTRYPOINT ["./root/run.sh"]