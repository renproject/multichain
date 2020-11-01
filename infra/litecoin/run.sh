#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/litecoind -conf=/root/.litecoin/litecoin.conf
sleep 10

# Print setup
echo "LITECOIN_ADDRESS=$ADDRESS"

# Import the address
/app/bin/litecoin-cli importaddress $ADDRESS

# Import the private key to spend UTXOs
/app/bin/litecoin-cli importprivkey $PRIV_KEY

# Generate enough block to pass the maturation time
/app/bin/litecoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    # generate new ltc to the address
    /app/bin/litecoin-cli generatetoaddress 1 $ADDRESS
    sleep 5
    # send tx to own address while paying fee to the miner
    /app/bin/litecoin-cli sendtoaddress $ADDRESS 0.5 "" "" true
    sleep 5
done
