#!/bin/bash
MNEMONIC=$1
ADDRESS=$2

ganache-cli      \
  -h 0.0.0.0     \
  -a 105         \
  -k muirGlacier \
  -i 420         \
  -m "$MNEMONIC" \
  -p 8575        \
  -u $ADDRESS    \
  -b 1           \
  -l 60000000    \
  --chainId 420
