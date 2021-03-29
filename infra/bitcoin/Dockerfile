FROM ubuntu:bionic

RUN apt-get update && apt-get install --yes software-properties-common
RUN apt-get install --yes curl

# bitcoind PPA is unmaintained, so we download it from the bitcoincore website.
# This can sometimes stall, which can be fixed by killing and restarting the
# build.
RUN curl https://bitcoincore.org/bin/bitcoin-core-0.19.1/bitcoin-0.19.1-x86_64-linux-gnu.tar.gz | tar -xz
RUN mv ./bitcoin-0.19.1 /app
RUN chmod +x /app/bin/bitcoind
RUN chmod +x /app/bin/bitcoin-cli

COPY bitcoin.conf /root/.bitcoin/
COPY run.sh /root/
RUN chmod +x /root/run.sh

EXPOSE 18443

ENTRYPOINT ["./root/run.sh"]