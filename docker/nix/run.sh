#!/bin/bash
ADDRESS=$1

# Start
nixd -conf=/root/.nix/nix.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "NIX_ADDRESS=$ADDRESS"

# Import the address
nix-cli importaddress $ADDRESS

# Generate enough blocks to pass the maturation time
nix-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    nix-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
