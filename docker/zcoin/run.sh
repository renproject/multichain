#!/bin/bash
ADDRESS=$1

# Start
zcoind -conf=/root/.zcoin/zcoin.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "ZCOIN_ADDRESS=$ADDRESS"

# Import the address
zcoin-cli importaddress $ADDRESS

# Generate enough blocks to pass the maturation time
zcoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    zcoin-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
