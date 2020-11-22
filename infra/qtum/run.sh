#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/qtumd
sleep 10

# Print setup
echo "QTUM_ADDRESS=$ADDRESS"

# Import the address
/app/bin/qtum-cli importaddress $ADDRESS

# Import the private key to spend UTXOs
/app/bin/qtum-cli importprivkey $PRIV_KEY

# Generate enough block to pass the maturation time
/app/bin/qtum-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/qtum-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
