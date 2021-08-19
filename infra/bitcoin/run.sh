#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/bitcoind -regtest -daemon
sleep 20

# Print setup
echo "BITCOIN_ADDRESS=$ADDRESS"

# Create wallet
/app/bin/bitcoin-cli createwallet "testwallet"

# Import the address
/app/bin/bitcoin-cli -regtest importaddress $ADDRESS

# Import the private key to spend UTXOs
/app/bin/bitcoin-cli -regtest importprivkey $PRIV_KEY

# Generate enough block to pass the maturation time
/app/bin/bitcoin-cli -regtest generatetoaddress 101 $ADDRESS

# Simulate mining
while :
do
    # generate new btc to the address
    /app/bin/bitcoin-cli -regtest generatetoaddress 1 $ADDRESS
    sleep 5
    # send tx to own address while paying fee to the miner
    /app/bin/bitcoin-cli -regtest sendtoaddress $ADDRESS 1 "" "" true
    sleep 5
done
