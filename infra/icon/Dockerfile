# Build rewardcalculator
FROM golang:1.12-alpine as builder

ARG RC_REPO=icon-project/rewardcalculator

RUN apk add --no-cache git make curl
ADD https://api.github.com/repos/icon-project/rewardcalculator/git/refs/heads/master version.json

RUN RC_TAG=$(curl --silent "https://api.github.com/repos/$RC_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') \
&& git clone --branch $RC_TAG --single-branch https://github.com/icon-project/rewardcalculator.git /rewardcalculator

RUN cd /rewardcalculator && make

# Build T-Bears
FROM python:3.7-slim-stretch

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    git \
    libc-dev \
    pkg-config \
    rabbitmq-server \
    vim-tiny \
    lsof \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work
COPY --from=builder /rewardcalculator/bin/icon_rc /usr/local/bin/
COPY --from=builder /rewardcalculator/bin/rctool /usr/local/bin/
COPY entry.sh /usr/local/bin
COPY tbears_server_config.json .
ADD genesis genesis

RUN chmod +x ./genesis/install.sh
RUN chmod +x ./genesis/init.sh
RUN chmod +x /usr/local/bin/entry.sh

RUN ./genesis/install.sh \
    && rm -rf /root/.cache/pip
RUN ./genesis/init.sh \
    && rm -rf genesis

EXPOSE 9000
ENTRYPOINT ["/usr/local/bin/entry.sh"]