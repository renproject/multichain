#!/bin/bash
ADDRESS=$1

# Start
/app/bin/digibyted
sleep 10

# Print setup
echo "DIGIBYTE_ADDRESS=$ADDRESS"

# Import the address
/app/bin/digibyte-cli importaddress $ADDRESS

# Generate enough block to pass the maturation time
/app/bin/digibyte-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    /app/bin/digibyte-cli generatetoaddress 1 $ADDRESS
    sleep 10
done
