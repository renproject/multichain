#!/bin/bash

UNIX_TIMESTAMP=$(date +"%s000")

echo $UNIX_TIMESTAMP

/app/build/bin/linux-amd64/node -log-level DEBUG -min-peers-mining 0 -outdate 1000d -api-address 0.0.0.0:6869 -grpc-address 0.0.0.0:6870 -blockchain-type integration -integration.account-seed CApJrVsZ6AY5zbunL2nqgrb7MkJF9rPiFz63RtaRPyna -integration.genesis.timestamp $UNIX_TIMESTAMP -integration.genesis.block-timestamp $UNIX_TIMESTAMP -integration.genesis.signature 5TdjutbC2P1AFqXp5NvXgyjjGDuoS2R9BHXutkygsMnae6hNkEMBvTD4HPBVc1m9jxFzbe7xCKGkXVpWdC8R1qF1 -integration.address-scheme-character I -build-extended-api -serve-extended-api