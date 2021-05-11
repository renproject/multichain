#!/bin/bash

NODE_DIR=/root/.bor
DATA_DIR=/root/.bor/data

bor --datadir $DATA_DIR init $NODE_DIR/genesis.json
cp $NODE_DIR/nodekey $DATA_DIR/bor/
cp $NODE_DIR/static-nodes.json $DATA_DIR/bor/

touch /root/logs/bor.log
bor --datadir $DATA_DIR \
  --port 30303 \
  --bor.heimdall http://heimdall$INDEX:1317 \
  --http --http.addr '0.0.0.0' \
  --http.vhosts '*' \
  --http.corsdomain '*' \
  --http.port 8545 \
  --http.api 'personal,eth,net,web3,txpool,miner,admin,bor' \
  --syncmode 'full' \
  --networkid '15001' \
  --miner.gaslimit '2000000000' \
  --txpool.nolocals \
  --txpool.accountslots '128' \
  --txpool.globalslots '20000' \
  --txpool.lifetime '0h16m0s' \
  --unlock $ADDRESS \
  --keystore $NODE_DIR/keystore \
  --password $NODE_DIR/password.txt \
  --allow-insecure-unlock \
  --mine > /root/logs/bor.log 2>&1 &
tail -f /root/logs/bor.log
