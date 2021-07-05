#!/bin/bash
MNEMONIC=$1
ADDRESS=$2

ganache-cli      \
  -h 0.0.0.0     \
  -a 105         \
  -k muirGlacier \
  -l 15000000    \
  -i 420         \
  -b 1           \
  -m "$MNEMONIC" \
  -u $ADDRESS    \
  --chainId 1337
