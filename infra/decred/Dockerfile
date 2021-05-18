FROM ubuntu:xenial

ENV VERSION v1.6.2
ENV RELEASE_NAME decred-linux-amd64-${VERSION}
ENV RELEASE_FILE ${RELEASE_NAME}.tar.gz
ENV DCR_URL https://github.com/decred/decred-binaries/releases/download/${VERSION}/${RELEASE_FILE}


RUN apt-get update && apt-get install --yes software-properties-common
RUN apt-get install --yes curl

RUN cd /tmp
RUN curl -SLO ${DCR_URL}
Run tar xvzf ${RELEASE_FILE}
Run mkdir /decred 
RUN mv ${RELEASE_NAME}/dcrd /decred
RUN mv ${RELEASE_NAME}/dcrwallet /decred
RUN mv ${RELEASE_NAME}/dcrctl /decred
RUN chmod +x /decred/dcrd
RUN chmod +x /decred/dcrwallet
RUN chmod +x /decred/dcrctl

COPY dcrd.conf /root/.dcrd/
COPY rpc.key /root/.dcrd/
COPY rpc.cert /root/.dcrd/
COPY run.sh /root/
RUN chmod +x /root/run.sh

# Cleanup
RUN cd /root/ && rm -rf /tmp/*

# PEER & RPC PORTS
EXPOSE 19556 18555

ENTRYPOINT ["./root/run.sh"]
