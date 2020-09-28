#!/bin/bash
ADDRESS=$1

# Print setup
echo "TERRA_ADDRESS=$ADDRESS"

# Register client key
terracli keys add validator --keyring-backend=test
echo $(terracli keys show validator)

# Initialize tesnet
terrad init testnet --chain-id testnet
terrad add-genesis-account $(terracli keys show validator -a --keyring-backend=test) 10000000000uluna
terrad add-genesis-account $ADDRESS 10000000000uluna,10000000000ukrw,10000000000uusd,10000000000usdr,10000000000umnt
terrad gentx --amount 10000000000uluna --name validator --keyring-backend=test
terrad collect-gentxs

# Start terrad
terrad start --rpc.laddr "tcp://0.0.0.0:26657"
