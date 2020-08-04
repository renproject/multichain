#!/bin/bash
WIF=$1

# Create wallet from WIF
echo $WIF > source
echo -e "1\n1\n" | /var/ontology/ontology account import -s source --wif
echo "1" | /var/ontology/ontology --testmode --gasprice 0
sleep 15
