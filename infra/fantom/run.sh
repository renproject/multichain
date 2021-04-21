#!/bin/bash
KEY=$1

lachesis               \
  --fakenet "1/1"      \
  --rpc                \
  --rpcvhosts "*"      \
  --rpcaddr "0.0.0.0"
  --nodekeyhex "$KEY"