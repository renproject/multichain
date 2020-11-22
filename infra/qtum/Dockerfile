# DEZU: Originally stolen from qtum-docker github page
FROM ubuntu
# MAINTAINER pandazwb

RUN apt-get update && apt-get install --yes -qq --no-install-recommends ca-certificates curl wget apt-utils jq
RUN apt-get install --yes curl

# install qtum binaries
RUN wget -c https://github.com/qtumproject/qtum/releases/download/mainnet-ignition-v0.20.1/qtum-0.20.1-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./qtum-0.20.1 /app
RUN chmod +x /app/bin/qtumd
RUN chmod +x /app/bin/qtum-cli

COPY qtum.conf /root/.qtum/
COPY run.sh /root/
RUN chmod +x /root/run.sh

EXPOSE 3888 3889 13888 13889

ENTRYPOINT ["./root/run.sh"]
