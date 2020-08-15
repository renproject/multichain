FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget

RUN wget https://github.com/Jenova7/dev_purpose/releases/download/2.1.3T/electra-2.1.3-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./electra-2.1.3 /app

RUN chmod +x /app/bin/electrad
RUN chmod +x /app/bin/electra-cli

COPY Electra.conf /root/.electra/
COPY run.sh /root/
RUN chmod +x /root/run.sh

# Regtest network RPC port
EXPOSE 15788

ENTRYPOINT ["./root/run.sh"]