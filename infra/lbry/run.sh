#!/bin/bash
ADDRESS=$1

# Start
/app/lbrycrdd -conf=/root/.lbrycrd/lbrycrd.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "LBRY_ADDRESS=$ADDRESS"

# Import the address
/app/lbrycrd-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/lbrycrd-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/lbrycrd-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
