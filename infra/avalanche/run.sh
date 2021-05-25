#!/bin/bash

AVAX_USERNAME=$1
AVAX_PASSWORD=$2
AVAX_PK=$3
AVAX_ADDRESS=$4
C_AVAX_PK=$5
C_AVAX_HEX_ADDRESS=$6
C_AVAX_BECH32_ADDRESS=$7

avalanchego \
	--assertions-enabled=true \
	--tx-fee=1000000 \
	--public-ip=0.0.0.0 \
	--network-id=local \
	--signature-verification-enabled=true \
	--api-admin-enabled=true \
	--api-ipcs-enabled=true \
	--api-keystore-enabled=true \
	--api-metrics-enabled=true \
	--http-host=0.0.0.0 \
	--http-port=9650 \
	--http-tls-enabled=false \
	--plugin-dir=/app/avalanchego-v1.4.5/avalanchego-latest/plugins \
	--log-level=info \
	--snow-avalanche-batch-size=30 \
	--snow-avalanche-num-parents=5 \
	--snow-sample-size=1 \
	--snow-quorum-size=1 \
	--snow-virtuous-commit-threshold=1 \
	--snow-rogue-commit-threshold=4 \
	--staking-enabled=false \
	--staking-port=9651 \
	--api-auth-required=false &

# create a new user
sleep 10
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"keystore.createUser",
    "params" :{
        "username":"'"$AVAX_USERNAME"'",
	      "password":"'"$AVAX_PASSWORD"'"
    }
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/keystore

# import private key that contains AVAX into X-chain
sleep 1
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"avm.importKey",
    "params" :{
        "username":"'"$AVAX_USERNAME"'",
	      "password":"'"$AVAX_PASSWORD"'",
	      "privateKey":"'"$AVAX_PK"'"
    }
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/bc/X

# import private key into C-chain
sleep 1
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"avax.importKey",
    "params" :{
        "username" :"'"$AVAX_USERNAME"'",
        "password":"'"$AVAX_PASSWORD"'",
        "privateKey":"'"$AVAX_PK"'"
    }
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/bc/C/avax

# export the AVAX to C-chain
sleep 1
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"avm.exportAVAX",
    "params" :{
        "to":"'"$C_AVAX_BECH32_ADDRESS"'",
        "destinationChain": "C",
        "amount": 5000000000000,
        "username":"'"$AVAX_USERNAME"'",
        "password":"'"$AVAX_PASSWORD"'"
    }
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/bc/X

# import AVAX to the hex address
sleep 1
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"avax.importAVAX",
    "params" :{
        "to":"'"$C_AVAX_HEX_ADDRESS"'",
        "sourceChain":"X",
        "username":"'"$AVAX_USERNAME"'",
        "password":"'"$AVAX_PASSWORD"'"
    }
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/bc/C/avax

wait %1
