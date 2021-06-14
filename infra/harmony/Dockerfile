FROM golang:1.14

WORKDIR $GOPATH/src/github.com/harmony-one

RUN apt update && apt upgrade -y && apt update -y && apt install unzip \
libgmp-dev libssl-dev curl git jq make gcc g++ bash sudo -y

RUN git clone https://github.com/harmony-one/harmony.git
RUN git clone https://github.com/harmony-one/bls.git
RUN git clone https://github.com/harmony-one/mcl.git
RUN cd harmony && make

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 9500
EXPOSE 9501
EXPOSE 9599
EXPOSE 9598

WORKDIR /

ENTRYPOINT ["./root/run.sh"]