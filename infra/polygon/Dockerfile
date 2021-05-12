FROM maticnetwork/bor:master

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

RUN mkdir -p /root/.bor/keystore
COPY genesis.json /root/.bor/genesis.json
COPY nodekey /root/.bor/nodekey
COPY static-nodes.json /root/.bor/static-nodes.json
COPY json-keystore /root/.bor/keystore/UTC--2021-05-11T14-27-08.753Z--0xbf7A416377ed8f1F745A739C8ff59094EB2FEFD2
COPY password.txt /root/.bor/password.txt

ENTRYPOINT [ "./root/run.sh" ]
