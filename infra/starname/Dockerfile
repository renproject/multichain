FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget
RUN wget -c https://github.com/iov-one/iovns/releases/download/v0.9.4/iovns-0.9.4-linux-amd64.tar.gz -O - | tar xz
RUN mkdir -p /app/bin
RUN mv ./iovnsd ./iovnscli /app/bin
RUN chmod +x /app/bin/iovnsd
RUN chmod +x /app/bin/iovnscli
RUN ln -s /app/bin/iovnsd /usr/bin/iovnsd
RUN ln -s /app/bin/iovnscli /usr/bin/iovnscli

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 46657

ENTRYPOINT ["./root/run.sh"]
