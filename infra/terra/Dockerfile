FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common build-essential wget curl git

RUN wget -c https://golang.org/dl/go1.16.8.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.8.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app
RUN git clone https://github.com/terra-money/classic-core.git
WORKDIR /app/classic-core
RUN git fetch --all -p
RUN git checkout v0.5.5
RUN make install

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 26657

ENV PATH=$PATH:/root/go/bin

WORKDIR /

ENTRYPOINT ["./root/run.sh"]
