#!/bin/bash
ADDRESS=$1

# Start
/app/bin/dashd -conf=/root/.dashcore/dash.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10

# Print setup
echo "DASH_ADDRESS=$ADDRESS"

# Import the address
/app/bin/dash-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/dash-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/dash-cli generatetoaddress 1 $ADDRESS
    sleep 10
done