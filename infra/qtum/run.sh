#!/bin/bash
ADDRESS=$1
PRIV_KEY=$2

# Start
/app/bin/qtumd
sleep 10

# Print setup
echo "(QTUM): QTUM_ADDRESS=$ADDRESS"

# Import the address
/app/bin/qtum-cli importaddress $ADDRESS

# Import the private key to spend UTXOs
/app/bin/qtum-cli importprivkey $PRIV_KEY

echo "(QTUM): Trying to generate 501 blocks..." # DEZU: TODO: Remove debug
# Generate enough block to pass the maturation time (500 for Qtum)
/app/bin/qtum-cli generatetoaddress 501 $ADDRESS
echo "(QTUM): Blocks (hopefully) generated!" # DEZU: TODO: Remove debug

#echo "(QTUM): Running 'getbalance'..." # DEZU: TODO: Remove debug
#/app/bin/qtum-cli getbalance

#echo "(QTUM): Running 'getwalletinfo'..." # DEZU: TODO: Remove debug
#/app/bin/qtum-cli getwalletinfo

#echo "(QTUM): Running 'estimatesmartfee 10'..." # DEZU: TODO: Remove debug
#/app/bin/qtum-cli estimatesmartfee 10

#echo "(QTUM): Running 'listunspent 0 999999999 qb15NCu3w4zyd14L21P99AdqmovHCiCEqC'..." # DEZU: TODO: Remove debug
#/app/bin/qtum-cli listunspent 0 999999999 qb15NCu3w4zyd14L21P99AdqmovHCiCEqC | cat

# Simulate mining
while :
do
    echo "(QTUM): Running 'generatetoaddress 1'..." # DEZU: TODO: Remove debug
    /app/bin/qtum-cli generatetoaddress 1 $ADDRESS
    sleep 5
    # send tx to own address while paying fee to the miner
    echo "(QTUM): Running 'sendtoaddress $ADDRESS 0.5 "" "" true'..." # DEZU: TODO: Remove debug
    /app/bin/qtum-cli sendtoaddress $ADDRESS 0.5 "" "" true
    sleep 5
    #echo "(QTUM): Running 'estimatesmartfee 15'..." # DEZU: TODO: Remove debug
    #/app/bin/qtum-cli estimatesmartfee 10
done
