#!/bin/bash
ADDRESS=$1

# Start
/app/bin/peercoind
sleep 10

# Print setup
echo "BITCOIN_ADDRESS=$ADDRESS"

# Import the address
/app/bin/peercoin-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/peercoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/peercoin-cli generatetoaddress 1 $ADDRESS
    sleep 10
done