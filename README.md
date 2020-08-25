# `🔗 multichain`

## Layout

`/` declares the assets and chains that exist, but provides no chain-specific implementations.

`/infra` defines a local deployment of the multichain using `docker-compose`. All underlying chains provide a `Dockerfile` and service definition to make running node instances easy. All chains need to add a `Dockerfile` and service definition that allows the multichain to spin up a local development-mode instance of the chain. This is necessary for running comprehensive local test suites.

`/api` defines the different compatibility APIs that exist: Account, Address, Contract, Gas, and UTXO. Chains should implement the APIs that are relevant to them. For example, Bitcoin (and its forks) implements the Address, Gas, and UTXO APIs. No actual implementations should be added to this folder.

`/chain`  defines all of the chain-specific implementations of the APIs. Each chain has its own sub-package. For example, Bitcoin, Bitcoin Cash, Dogecoin, and Zcash are all chains that implement the Address, Gas, and UTXO APIs, and each of these implementations are in `/chain/bitcoin`, `/chain/bitcoincash`, `/chain/dogecoin`, and `/chain/zcash` respectively.

## Example

The `🔗 multichain` is designed to be flexible enough to support any kind of chain. Anyone is free to contribute to the `🔗 multichain` by adding support for a new chain, or improving support for an existing chain. To show how this is done, we will walk-through an example: adding support for Dogecoin.

### Chains and Assets

Before doing anything else, let's add an enumeration for the `Asset` and `Chain` types, which can be found in `package multichain`. To avoid favouritism, all assets and chains are listed in alphabetical order. Unless otherwise advised by an offiical team member, the names and tickers found on https://coinmarketcap.com must be used.

Adding an `Asset`:

```go
// Enumeration of supported assets. When introducing a new chain, or new asset
// from an existing chain, you must add a human-readable string to this set of
// enumerated values. Assets must be listed in alphabetical order.
const (
    BCH  = Asset("BCH")  // Bitcoin Cash
    BTC  = Asset("BTC")  // Bitcoin
    DOGE = Asset("DOGE") // Dogecoin (This is our new asset!)
    ETH  = Asset("ETH")  // Ether
    ZEC  = Asset("ZEC")  // Zcash
)
```

Adding a `Chain`:

```go
// Enumeration of supported chains. When introducing a new chain, you must add a
// human-readable string to this set of enumerated values. Chains must be listed
// in alphabetical order.
const (
    Acala       = Chain("Acala")
    Bitcoin     = Chain("Bitcoin")
    BitcoinCash = Chain("BitcoinCash")
    Dogecoin    = Chain("Dogecoin") // (This is our new chain!)
    Ethereum    = Chain("Ethereum")
    Zcash       = Chain("Zcash")
)
```

### Docker

Next, we need to setup a Docker container in the `/infra` folder. This is needed for local test suites, allowing for end-to-end integrated testing directly against a node. Doing this requires a couple of steps.

First, we create a new `dogecoin/` folder in the `/infra` folder:

```
/infra
|-- /bitcoin
|-- /bitcoincash
|-- /dogecoin         # This is our new folder!
|   |-- Dockerfile    # This is our new Dockerfile!
|   |-- dogecoin.conf
|   |-- run.sh        # This is our new run file!
|-- /zcash
|-- .env
|-- docker-compose.yaml
```

The new folder _must_ at least contain a `Dockerfile` that installs the node, and a `run.sh` file that runs the nodes. The node _should_ be run in test mode. The new folder can also contain other files that are specific to the needs of the chain being added. In our case, the `dogecoin.conf` file is also needed to configure the node. (We will omit showing all the code here, since there is quite a bit of it, but you can check it out in the `/infra/dogecoin` folder.)

Second, we add an entry to the `.env` file. Our entry _must_ include a private key that will have access to funds, and the public address associated with that private key. We will add:

```sh
#
# Dogecoin
#

# Address that will receive mining rewards. Generally, this is set to an address
# for which the private key is known by a test suite. This allows the test suite
# access to plenty of testing funds.
export DOGECOIN_PK=cRZnRgH2ztcJupCzkWbq2mjiT8PSFAmtYRYb1phg1vSRRcNBX4w4
export DOGECOIN_ADDRESS=n3PSSpR4zqUKWH4tcRjP9aTwJ4GmixQXmt
```

Last, we add a service to the `docker-compose.yaml` file. This allows the node to boot alongside the other nodes in the multichain. This entry must expose the node for use in tests, and must not overlap with other nodes that already exist (ports are reserved on a first-come-first-serve basis). We will define the service as:

```yaml
##
## Dogecoin
##

dogecoin:
  build:
    context: ./dogecoin
  ports:
    - "0.0.0.0:18332:18332"
  entrypoint:
    - "./root/run.sh"
    - "${DOGECOIN_ADDRESS}"
```

### Address API

All chains _should_ implement the Address API. Luckily for Dogecoin, it is so similar to Bitcoin that we can re-export the Bitcoin implementation without the need for custom modifications. In `/chain/dogecoin/address.go` we add:

```go
package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type (
	AddressEncoder       = bitcoin.AddressEncoder
	AddressDecoder       = bitcoin.AddressDecoder
	AddressEncodeDecoder = bitcoin.AddressEncodeDecoder
)
```

These three interfaces allow users of the `🔗 multichain` to easily encode and decode Dogecoin addresses. Other chains will need to provide their own implementations, based on their chains address standards.

### Gas API

Most, but not all, chains _should_ implement the Gas API. Again, Dogecoin is so similar to Bitcoin that we can re-export the Bitcoin implementation in `/chain/dogecoin/gas.go`:

```go
package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type GasEstimator = bitcoin.GasEstimator

var NewGasEstimator = bitcoin.NewGasEstimator
```

The interface allows users of the `🔗 multichain` to estimate gas prices (although, the current implementation is _very_ simple). The associated function allows users to construct an instance of the interface for Dogecoin.

### UTXO API

Generally speaking, chains fall into two categories: account-based or UTXO-based (and some can even be both). Bitcoin, and its forks, are all UTXO-based chains. As a fork of Bitcoin, Dogecoin is a UTXO-based chain, so we implement the UTXO API. To implement the UTXO API, we must implement the `Tx`, `TxBuilder`, and `Client` interfaces. More information can be found in the comments of `/api/utxo` folder.

Again, the implementation for Dogecoin is trivial. In `/chain/dogecoin/utxo`, we have:

```go
package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type (
	Tx            = bitcoin.Tx
	TxBuilder     = bitcoin.TxBuilder
	Client        = bitcoin.Client
	ClientOptions = bitcoin.ClientOptions
)

var (
	NewTxBuilder         = bitcoin.NewTxBuilder
	NewClient            = bitcoin.NewClient
	DefaultClientOptions = bitcoin.DefaultClientOptions
)
```

Up to this point, we have done nothing but re-export Bitcoin. So what makes Dogecoin different? And how can we express that difference? Well, the `/chain/dogecoin` folder is the place where we must define anything else Dogecoin users will need. In the case of Dogecoin, the only thing that differentiates it from Bitcoin is the `*chaincfg.Param` object. We define this in `/chain/dogecoin/dogecoin.go`:

```go
package dogecoin

import (
	"github.com/btcsuite/btcd/chaincfg"
)

func init() {
	if err := chaincfg.Register(&MainNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&RegressionNetParams); err != nil {
		panic(err)
	}
}

var MainNetParams = chaincfg.Params{
	Name: "mainnet",
	Net:  0xc0c0c0c0,

	// Address encoding magics
	PubKeyHashAddrID: 30,
	ScriptHashAddrID: 22,
	PrivateKeyID:     158,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x02, 0xfa, 0xc3, 0x98}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x02, 0xfa, 0xca, 0xfd}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dogecoin does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "doge",
}

var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// Dogecoin has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(dogecoin)), so that we can
	// register the regtest network.
	Net: 0xd063c017,

	// Address encoding magics
	PubKeyHashAddrID: 111,
	ScriptHashAddrID: 196,
	PrivateKeyID:     239,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dogecoin does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "dogert",
}
```

Most of the functions that we have re-exported expected `*chaincfg.Params` as an argument. By defining one for regnet and mainnet, users can construct Dogecoin instances of the UTXO API by using these params.

## Test Suite

1. Install Docker
2. Install Docker Compose
3. Run Docker
4. Run `./test.sh`

Example output:

```sh
Creating network "docker_default" with the default driver
Building bitcoin

...

Successfully built 1ebb03faa04f
Successfully tagged docker_bitcoin:latest
Building bitcoincash

...

Successfully built e12e98011869
Successfully tagged docker_bitcoincash:latest
Building zcash

...

Successfully built 56231a29ca2e
Successfully tagged docker_zcash:latest
docker_bitcoin_1 is up-to-date
docker_bitcoincash_1 is up-to-date
docker_zcash_1 is up-to-date
Waiting for multichain to boot...
=== RUN   TestMultichain
Running Suite: Multichain Suite
===============================

...

Stopping docker_bitcoincash_1 ... done
Stopping docker_zcash_1       ... done
Stopping docker_bitcoin_1     ... done
Removing docker_bitcoincash_1 ... done
Removing docker_zcash_1       ... done
Removing docker_bitcoin_1     ... done
Removing network docker_default
Done!
```
