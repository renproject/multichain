FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget

RUN wget -c https://github.com/dogecoin/dogecoin/releases/download/v1.14.2/dogecoin-1.14.2-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./dogecoin-1.14.2 /app

COPY dogecoin.conf /root/.dogecoin/dogecoin.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 18332

ENTRYPOINT ["./root/run.sh"]
