#!/bin/bash
ADDRESS=$1

# Start
electrad -conf=/root/.electra/Electra.conf
sleep 10

# Print setup
echo "ELECTRA_ADDRESS=$ADDRESS"

# Import the address
/app/bin/electra-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/electra-cli setgenerate true 101

# Simulate mining
while :
do
    /app/bin/electra-cli setgenerate true 1
    sleep 10
done