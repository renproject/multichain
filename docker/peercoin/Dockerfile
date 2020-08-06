FROM debian:10

RUN apt-get update
RUN apt-get install -y apt-transport-https wget gnupg

RUN sh -c "echo 'deb https://peercoin.github.io/deb-repo/ buster main' >> /etc/apt/sources.list.d/peercoin.list" && \
        wget -O - https://peercoin.github.io/deb-repo/peercoin.apt.key | apt-key add -

RUN apt-get update && \
        apt-get install -y peercoind peercoin-tx

COPY peercoin.conf /root/.peercoin/peercoin.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 9901 9902 9903 9904

ENTRYPOINT ["./root/run.sh"]
