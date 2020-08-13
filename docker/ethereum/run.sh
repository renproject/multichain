#!/bin/bash
MNEMONIC=$1
ADDRESS=$2

ganache-cli                     \
  -h 0.0.0.0                    \
  -f https://cloudflare-eth.com \
  -i 420                        \
  -m $MNEMONIC                  \
  -u $ADDRESS