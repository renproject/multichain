#!/bin/bash
ADDRESS=$1

# Start
/app/bin/bitcoind
sleep 10

# Print setup
echo "BITCOIN_ADDRESS=$ADDRESS"

# Import the address
/app/bin/bitcoin-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/bitcoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/bitcoin-cli generatetoaddress 1 $ADDRESS
    sleep 10
done