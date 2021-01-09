FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget

RUN wget -c https://github.com/floblockchain/flo/releases/download/v0.15.2.1/flo-0.15.2-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./flo-0.15.2 /app
RUN chmod +x /app/bin/flod
RUN chmod +x /app/bin/flo-cli

COPY flo.conf /root/.flo/flo.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 8333

ENTRYPOINT ["./root/run.sh"]
