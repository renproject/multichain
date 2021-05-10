#!/bin/bash
ADDRESS=$1

# Start
BGLd -conf=/root/.BGL/BGL.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "BGL_ADDRESS=$ADDRESS"

# Import the address
BGL-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time (100 blocks)
BGL-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    BGL-cli generatetoaddress 1 $ADDRESS
    sleep 10
done