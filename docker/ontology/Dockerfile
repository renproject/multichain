FROM ubuntu:xenial

RUN apt-get update && apt-get install --yes software-properties-common
RUN apt-get install --yes wget

RUN wget https://github.com/ontio/ontology/releases/download/v2.0.0/ontology-linux-amd64

ENV ONTOLOGY_PATH /var/ontology
RUN mkdir -p $ONTOLOGY_PATH
RUN mv ./ontology-linux-amd64 $ONTOLOGY_PATH/ontology
RUN chmod +x $ONTOLOGY_PATH/ontology
COPY run.sh $ONTOLOGY_PATH
RUN chmod +x $ONTOLOGY_PATH/run.sh

EXPOSE 20336

ENTRYPOINT ["./var/ontology/run.sh"]