#!/bin/bash
MNEMONIC=$1
ADDRESS=$2

ganache-cli      \
  -h 0.0.0.0     \
  -a 105         \
  -k muirGlacier \
  -i 421         \
  -m "$MNEMONIC" \
  -p 8565        \
  -u $ADDRESS    \
  -b 1           \
  -l 60000000    \
  --chainId 421
