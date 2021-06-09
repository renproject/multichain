#!/bin/bash
ADDRESS=SsaGEEZu2L8x93qvKzzahtzQ7yzkec3i8wL
PRIV_KEY=PsUQEpYDXVwphd9xNXUMj63LyxSWPTor3RDgfw9DMdH9tDkJaosyp

# Start
screen -dm /decred/dcrd --simnet --miningaddr=$ADDRESS
screen -dm /decred/dcrwallet --simnet --createtemp
sleep 2

# Print setup
echo "DECRED_ADDRESS=$ADDRESS"

/decred/dcrctl --simnet generate 101

# Simulate mining
while :
do
	/decred/dcrctl --simnet generate 10
    sleep 5
done

