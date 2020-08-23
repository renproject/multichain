FROM debian:buster

RUN apt-get update --fix-missing
RUN apt-get install --yes curl bzip2 libssl-dev

RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get install --yes nodejs

RUN curl -sSf https://raw.githubusercontent.com/solana-labs/solana/v1.2.13/install/solana-install-init.sh | sh -s - v1.2.13

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 8899/tcp
EXPOSE 8900
EXPOSE 9900

ENTRYPOINT ["./root/run.sh"]
CMD [""]