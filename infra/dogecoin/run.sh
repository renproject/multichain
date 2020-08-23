#!/bin/bash
ADDRESS=$1

# Start
/app/bin/dogecoind -conf=/root/.dogecoin/dogecoin.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "DOGECOIN_ADDRESS=$ADDRESS"

# Import the address
/app/bin/dogecoin-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/dogecoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/dogecoin-cli generatetoaddress 1 $ADDRESS
    sleep 10
done