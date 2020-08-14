FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget
RUN wget -c https://github.com/terra-project/core/releases/download/v0.4.0-rc.2/terra_0.4.0_rc.2_Linux_x86_64.tar.gz -O - | tar xz
RUN mkdir -p /app/bin
RUN mv ./terrad ./terracli /app/bin
RUN chmod +x /app/bin/terrad
RUN chmod +x /app/bin/terracli
RUN ln -s /app/bin/terrad /usr/bin/terrad
RUN ln -s /app/bin/terracli /usr/bin/terracli

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 26657

ENTRYPOINT ["./root/run.sh"]
