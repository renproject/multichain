#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/bitcoind
sleep 10

# Print setup
echo "BITCOIN_ADDRESS=$ADDRESS"

# Import the address
/app/bin/bitcoin-cli importaddress $ADDRESS

# Import the private key to spend UTXOs
/app/bin/bitcoin-cli importprivkey $PRIV_KEY

# Generate enough block to pass the maturation time
/app/bin/bitcoin-cli generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    # generate new btc to the address
    /app/bin/bitcoin-cli generatetoaddress 1 $ADDRESS
    sleep 5
    # send tx to own address while paying fee to the miner
    /app/bin/bitcoin-cli sendtoaddress $ADDRESS 0.5 "" "" true
    sleep 5
done
