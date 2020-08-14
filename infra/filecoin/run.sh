#!/bin/bash

./lotus daemon --lotus-make-genesis=dev.gen --genesis-template=localnet.json --bootstrap=false &

PID=$!

sleep 10

./lotus wallet import ~/.genesis-sectors/pre-seal-t01000.key

./lotus wallet import /root/miner.key

kill $PID

echo '
# Default config:
[API]
ListenAddress = "/ip4/0.0.0.0/tcp/1234/http"
RemoteListenAddress = "127.0.0.1:1234"
Timeout = "30s"
#
[Libp2p]
#  ListenAddresses = ["/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"]
#  AnnounceAddresses = []
#  NoAnnounceAddresses = []
#  ConnMgrLow = 150
#  ConnMgrHigh = 180
#  ConnMgrGrace = "20s"
#
[Pubsub]
#  Bootstrapper = false
#  RemoteTracer = "/ip4/147.75.67.199/tcp/4001/p2p/QmTd6UvR47vUidRNZ1ZKXHrAFhqTJAD27rKL9XYghEKgKX"
#
[Client]
#  UseIpfs = false
#  IpfsMAddr = ""
#  IpfsUseForRetrieval = false
#
[Metrics]
#  Nickname = ""
#  HeadNotifs = false
#' > ~/.lotus/config.toml

./lotus daemon --lotus-make-genesis=/root/dev.gen --genesis-template=/root/localnet.json --bootstrap=false &

sleep 5

./lotus-storage-miner init --genesis-miner --actor=t01000 --sector-size=2KiB --pre-sealed-sectors=~/.genesis-sectors --pre-sealed-metadata=~/.genesis-sectors/pre-seal-t01000.json --nosync

./lotus-storage-miner run --nosync &

sleep 15

MAIN_WALLET="$(jq -r '.t01000.Owner' ~/.genesis-sectors/pre-seal-t01000.json)"

./lotus send --from $MAIN_WALLET t17lz42vyixtlfs4a3j76cw53w32qb57w7g4g6qua 1000000

while :
do
    sleep 10
done
