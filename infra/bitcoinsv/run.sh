#!/bin/bash
ADDRESS=$1

# Start
bitcoind
sleep 10

# Print setup
echo "BITCOINSV_ADDRESS=$ADDRESS"

# Import the address
bitcoin-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
bitcoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    bitcoin-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
