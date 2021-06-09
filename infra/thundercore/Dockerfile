FROM ubuntu:xenial

RUN apt-get update --fix-missing
RUN apt-get install --yes curl

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install --yes nodejs git
RUN useradd -ms /bin/bash node

# workaround for npm not running prepare script when running as root
USER node
RUN npm config set prefix /home/node
RUN npm install -g thundercore/ganache-cli#tt
USER root
RUN ln -s /home/node/bin/ganache-cli /usr/bin/ganache-cli

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 8565

ENTRYPOINT ["./root/run.sh"]
