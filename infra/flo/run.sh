#!/bin/bash
ADDRESS=$1

# Start
/app/bin/flod -conf=/root/.flo/flo.conf # -server -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0 -rpcuser=user -rpcpassword=password
sleep 10


# Print setup
echo "FLO_ADDRESS=$ADDRESS"

# Import the address
/app/bin/flo-cli importaddress $ADDRESS

echo "done"

# Generate enough block to pass the maturation time
/app/bin/flo-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/flo-cli generatetoaddress 1 $ADDRESS
    sleep 10
done