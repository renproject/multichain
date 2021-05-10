FROM ubuntu:xenial

RUN apt-get update --fix-missing && apt-get install --yes software-properties-common wget

RUN wget -c https://github.com/dashpay/dash/releases/download/v0.15.0.0/dashcore-0.15.0.0-x86_64-linux-gnu.tar.gz -O - | tar xz
RUN mv ./dashcore-0.15.0 /app
RUN chmod +x /app/bin/dashd
RUN chmod +x /app/bin/dash-cli

COPY dash.conf /root/.dashcore/dash.conf
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

EXPOSE 21443

ENTRYPOINT ["./root/run.sh"]
