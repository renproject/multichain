FROM golang:1.14

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common git
RUN git clone https://github.com/wavesplatform/gowaves.git /app
WORKDIR /app
RUN ls -la
RUN make release-node

COPY run.sh ./root/run.sh
RUN chmod +x ./root/run.sh

EXPOSE 6869
EXPOSE 6870

ENTRYPOINT ["./root/run.sh"]
