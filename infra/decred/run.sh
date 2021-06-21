#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
screen -dm /decred/dcrd --simnet --miningaddr=$ADDRESS
screen -dm /decred/dcrwallet --simnet --createtemp
sleep 10

# Print setup
echo "DECRED_ADDRESS=$ADDRESS"

/decred/dcrctl --simnet --wallet importprivkey $PRIV_KEY

/decred/dcrctl --simnet generate 101

# Simulate mining
while :
do
	/decred/dcrctl --simnet generate 10
    sleep 5
done

