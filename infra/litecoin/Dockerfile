FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common curl

RUN curl https://download.litecoin.org/litecoin-0.18.1/linux/litecoin-0.18.1-x86_64-linux-gnu.tar.gz | tar -xz
RUN mv ./litecoin-0.18.1 /app
RUN chmod +x /app/bin/litecoind
RUN chmod +x /app/bin/litecoin-cli

COPY litecoin.conf /root/.litecoin/
COPY run.sh /root/
RUN chmod +x /root/run.sh

EXPOSE 19443

ENTRYPOINT ["./root/run.sh"]