FROM node:16-alpine

WORKDIR /root/app

COPY package.json .
RUN npm install

COPY hardhat.config.js .
COPY run.sh .
RUN chmod +x run.sh

EXPOSE 8545

ENTRYPOINT ["./run.sh"]
