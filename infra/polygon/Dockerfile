FROM ubuntu:groovy

# UPDATE
RUN apt update -y
RUN apt install -y build-essential git jq wget

# INSTALL GO
RUN wget -c https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

# INSTALL RABBITMQ
COPY install-rabbitmq.sh .
RUN chmod +x install-rabbitmq.sh
RUN /bin/bash install-rabbitmq.sh

# INSTALL NODEJS AND NPM
RUN apt-get install -y nodejs npm

# INSTALL DEPENDENCIES
RUN npm install -g solc@0.5.11
RUN npm install -g ganache-cli

# INSTALL MATIC-CLI
RUN npm install -g @maticnetwork/matic-cli

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 8545

ENTRYPOINT ["./root/run.sh"]
