FROM ubuntu:bionic

# Install zcashd.
RUN apt-get update                                                                                    && \
    apt-get install -y --no-install-recommends apt-transport-https gnupg2 ca-certificates wget        && \
    wget -qO - https://apt.z.cash/zcash.asc | apt-key add -                                           && \
    echo "deb [arch=amd64] https://apt.z.cash/ stretch main" | tee /etc/apt/sources.list.d/zcash.list && \
    apt-get update && apt-get install -y --no-install-recommends zcash                                && \
    mkdir -p /root/.zcash-params && zcash-fetch-params

COPY zcash.conf /root/.zcash/zcash.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 18232

ENTRYPOINT ["./root/run.sh"]
