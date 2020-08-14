#!/bin/bash
ADDRESS=$1

# Start
zcashd -mineraddress=$ADDRESS
sleep 10

echo "ZCASH_ADDRESS=$ADDRESS"

# Import the address
zcash-cli importaddress $ADDRESS

# Generate enough block to pass the maturation tim=
zcash-cli generate 101

# Simulate mining
while :
do
    zcash-cli generate 1
    sleep 10
done