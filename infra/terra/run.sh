#!/bin/bash
ADDRESS=$1

# Print setup
echo "TERRA_ADDRESS=$ADDRESS"

# Register client key
terrad keys add validator --keyring-backend=test
echo $(terrad keys show validator --keyring-backend=test)

# Initialize tesnet
terrad init testnet --chain-id testnet
terrad add-genesis-account $(terrad keys show validator -a --keyring-backend=test) 10000000000uluna
terrad add-genesis-account $ADDRESS 10000000000uluna,10000000000ukrw,10000000000uusd,10000000000usdr,10000000000umnt
terrad gentx validator 10000000000uluna --keyring-backend=test --chain-id=testnet
terrad collect-gentxs

# Start terrad
terrad start --rpc.laddr "tcp://0.0.0.0:26657"
