#!/bin/bash
ADDRESS=$1

# Print setup
echo "IOV_ADDRESS=$ADDRESS"

# Register client key
iovnscli keys add validator --keyring-backend=test
echo $(iovnscli keys show validator)

# Initialize testnet
iovnsd init testnet --chain-id testnet
iovnsd add-genesis-account $(iovnscli keys show validator -a --keyring-backend=test) 10000000000tiov
iovnsd add-genesis-account $ADDRESS 10000000000tiov
iovnsd gentx --amount 10000000000tiov --name validator --keyring-backend=test
iovnsd collect-gentxs
sed -i 's/stake/tiov/g' ~/.iovnsd/config/genesis.json
iovnsd validate-genesis

# Start iovnsd
iovnsd start --rpc.laddr "tcp://0.0.0.0:46657"
