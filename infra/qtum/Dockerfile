# DEZU: Modified version of file from qtum-docker github page
# I'm not 100% sure if all commands here are necessary, I've done a lot of copypasta
FROM ubuntu

RUN apt-get update --fix-missing && apt-get install --yes -qq --no-install-recommends ca-certificates curl wget apt-utils jq
RUN apt-get install --yes curl

# install qtum binaries
RUN wget -c https://github.com/qtumproject/qtum/releases/download/mainnet-ignition-v0.20.1/qtum-0.20.1-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./qtum-0.20.1 /app
RUN chmod +x /app/bin/qtumd
RUN chmod +x /app/bin/qtum-cli

COPY qtum.conf /root/.qtum/qtum.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

# These addresses are actually for testnet, but they work in regtest too apparently
EXPOSE 13888 13889

ENTRYPOINT ["./root/run.sh"]
