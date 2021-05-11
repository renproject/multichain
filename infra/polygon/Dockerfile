FROM maticnetwork/bor:master

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

COPY genesis.json /root/.bor/genesis.json
COPY nodekey /root/.bor/nodekey
COPY static-nodes.json /root/.bor/static-nodes.json

WORKDIR /
ENTRYPOINT [ "./root/run.sh" ]
