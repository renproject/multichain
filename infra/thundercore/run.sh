#!/bin/bash
MNEMONIC=$1
ADDRESS=$2

ganache-cli                     \
  -h 0.0.0.0                    \
  -a 105                        \
  -k pala                       \
  -i 420                        \
  -m "$MNEMONIC"                \
  -p 8565                       \
  -u $ADDRESS
