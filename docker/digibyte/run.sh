#!/bin/bash
ADDRESS=$1

# Start
digibyted
sleep 10

# Print setup
echo "DIGIBYTE_ADDRESS=$ADDRESS"

# Import the address
digibyte-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
digibyte-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
 digibyte-cli generatetoaddress 1 $ADDRESS
 sleep 10
done