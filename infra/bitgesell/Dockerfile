FROM ubuntu:bionic

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget

RUN wget -c https://github.com/wu-emma/bitgesell/releases/download/0.1.3/bitgesell_0.1.3_amd64.deb
RUN apt-get install -y ./bitgesell_0.1.3_amd64.deb

COPY BGL.conf /root/.BGL/BGL.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 18475

ENTRYPOINT ["./root/run.sh"]
