# `🔗 multichain`

## Layout

`/` defines all of the functions/types/constants that are common to all underlying chains.

`/compat` contains all of the Compat APIs. These APIs are the interfaces required to be implemented by underlying chains that are similar to each other. For example, the `/compat/bitcoincompat` folder defines the Bitcoin Compat API. The `Address`, `Client`, `Tx`, and `TxBuilder` interfaces must be implemented by all underlying chains that want to be compatible with the Bitcoin runtime. It also defines implementations that are likely to be common to these underlying chains (although, each underlying chain can override whatever it needs to).

`/chain`  contains all the implementations of the Compat APIs. Each chain has its own sub-package. For example, Bitcoin, Bitcoin Cash, Dogecoin, and Zcash are all underyling chains that implement the Bitcoin Compat API (defined in `/compat/bitcoincompat`), and each of these implementations are in `/chain/bitcoin`, `/chain/bitcoincash`, `/chain/dogecoin`, and `/chain/zcash` respectively.

`/docker` defines a local deployment of the multichain using `docker-compose`. All underlying chains provide a `Dockerfile` and service definition to make running node instances easy.

`/runtime` contains all of the Runtime modules. There is exactly one Runtime module for each Compat API, and you can think of a Runtime module as a way of accessing all of the Compat APIs in one place. This is the primary package used by users of the multichain (including the RenVM execution engine). If your chain implements one of the existing Compat APIs, you will *not* need to modify any of the Runtime modules.

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

Next, we need to setup a Docker container in the `docker/` folder. This is needed for local test suites, allowing for end-to-end integrated testing directly against a node. Doing this requires a couple of steps.

First, we create a new `dogecoin/` folder in the `docker/` folder:

```
docker/
|-- bitcoin/
|-- bitcoincash/
|-- dogecoin/         # This is our new folder!
|   |-- Dockerfile    # This is our new Dockerfile!
|   |-- dogecoin.conf
|   |-- run.sh        # This is our new run file!
|-- zcash/
|-- docker-compose.env
|-- docker-compose.yaml
```

The new folder _must_ at least contain a `Dockerfile` that installs the node, and a `run.sh` file that runs the nodes. The node _should_ be run in test mode. The new folder can also contain other files that are specific to the needs of the chain being added. In our case, the `dogecoin.conf` file is also needed to configure the node. (We will omit showing all the code here, since there is quite a bit of it, but you can check it out in the `docker/dogecoin/` folder.)

Second, we add an entry to the `docker-compose.env` file. Our entry _must_ include a private key that will have access to funds, and the public address associated with that private key. We will add:

```sh
#
# Dogecoin
#

# Address that will receive mining rewards. Generally, this is set to an address
# for which the private key is known by a test suite. This allows the test suite
# access to plenty of testing funds.
export DOGECOIN_PK=cUJCHRMSUwkcofsHjFWBELT3yEAejokdKhyTNv3DScodYWzztBae
export DOGECOIN_ADDRESS=mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4
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

### Runtime

The final thing that is required before the `🔗 multichain` supports our new chain is an integration into the runtime, defined in `package runtime`. The exact requirements for integration into the runtime vary from chain-to-chain. To make life easier, there is a set of common interfaces, known as the _Compat API_, that can be implemented by new chains. The Compat API is a set of interfaces, designed with the intention for multiple implementations to exist. For example, the Bitcoin Compat API is used by Bitcoin, Bitcoin Cash, and Zcash.

The Compat API is defined by `package compat` (and is used by the `Runtime` type in `package runtime`). All of the interfaces in `package bitcoincompat` belong to the Bitcoin Compat API, all of the interfaces in `package ethereumcompat` belong to the Ethereum Compat API, all of the interfaces in `package substratecompat` belong to the Substrate Compat API, and so on. Similarly, the `BitcoinXX`, `EthereumXXX`, and `SubstrateXXX` methods (defined by the `Runtime` type in `package runtime`) are all abstractions over the respective Compat APIs, but do not need to be modified!

Dogecoin is a fork of Bitcoin, so it is natural that we will support it by implementing the Bitcoin Compat API. Dogecoin is, in fact, so similar to Bitcoin that the implementation is trivial. All implementations belond in `package chain`, so we will create a new `package dogecoin` in that directory. Here, we create the `dogecoin.go` file and fill it with:

```go
package dogecoin

import (
    "github.com/renproject/multichain/chain/bitcoin"
    "github.com/renproject/multichain/compat/bitcoincompat"
)

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Dogecoin transactions.
func NewTxBuilder() bitcoincompat.TxBuilder {
    return bitcoin.NewTxBuilder()
}

// The Tx type is copied from Bitcoin.
type Tx = bitcoin.Tx
```

For a coin as simple as Dogecoin, nothing else is required! For more complex examples, you can checkout `package bitcoincash` and `package zcash` which need to define their own address and transaction formats.

### Custom Runtimes

Not all chains are as simple as the Dogecoin chain, and an existing Compat API may not be sufficient for your needs. In these scenarios, a little more work is required.

1. Define your own compat package (e.g. `package myawesomechaincompat`) in the `compat/` folder.
2. Define your own compat interfaces in your new compat package.
3. Define your own compat methods on the `Runtime` type in `package runtime`. You will always need a `MyAwesomeChainDecodeAddress(pack.String) (myawesoemcompat.Address, error)` method. If your blockchain is programmable, then defining methods for querying relevant events is usually sufficient. Otherwise, building/submitting transactions is probably going to be required.

If in doubt, get in touch with the Ren team at https://t.me/renproject and we will help you out!

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
