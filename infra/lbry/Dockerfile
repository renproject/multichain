FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget unzip

RUN wget -c https://github.com/lbryio/lbrycrd/releases/download/v0.17.4.6/lbrycrd-linux-1746.zip -O tmp.zip && unzip tmp.zip -d lbrycrd && rm tmp.zip
RUN mv ./lbrycrd /app
RUN chmod +x /app/lbrycrdd
RUN chmod +x /app/lbrycrd-cli

COPY lbrycrd.conf /root/.lbrycrd/lbrycrd.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 29245

ENTRYPOINT ["./root/run.sh"]
