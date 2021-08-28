#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/bitcoind -regtest -daemon
sleep 20

# Print setup
echo "BITCOIN_ADDRESS=$ADDRESS"

/app/bin/bitcoin-cli createwallet "testwallet"

sleep 10

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
    /app/bin/bitcoin-cli -regtest -named sendtoaddress address=$ADDRESS amount=0.1 subtractfeefromamount=false fee_rate=1
    sleep 5
done
