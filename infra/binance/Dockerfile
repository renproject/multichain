FROM ubuntu:xenial

RUN apt-get update --fix-missing
RUN apt-get install --yes curl

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install --yes nodejs
RUN npm install -g ganache-cli

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 8575

ENTRYPOINT ["./root/run.sh"]
