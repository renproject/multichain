#!/bin/bash
ADDRESS=$1

# Start
cd /app
make run
sleep 10

# Print setup
echo "ACALA_ADDRESS=$ADDRESS"
