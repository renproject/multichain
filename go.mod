module github.com/renproject/multichain

go 1.14

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/codahale/blake2 v0.0.0-20150924215134-8d10d0420cbf
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/drand/drand v1.0.3-0.20200714175734-29705eaf09d4 // indirect
	github.com/etcd-io/bbolt v1.3.3 // indirect
	github.com/ethereum/go-ethereum v1.9.19
	github.com/filecoin-project/go-address v0.0.3 // indirect
	github.com/filecoin-project/go-amt-ipld v0.0.0-20191205011053-79efc22d6cdc // indirect
	github.com/filecoin-project/go-amt-ipld/v2 v2.1.0 // indirect
	github.com/filecoin-project/go-bitfield v0.1.0 // indirect
	github.com/filecoin-project/go-data-transfer v0.5.0 // indirect
	github.com/filecoin-project/go-fil-markets v0.3.2 // indirect
	github.com/filecoin-project/lotus v0.4.1 // indirect
	github.com/filecoin-project/sector-storage v0.0.0-20200723200950-ed2e57dde6df // indirect
	github.com/filecoin-project/specs-actors v0.6.2-0.20200724193152-534b25bdca30 // indirect
	github.com/hannahhoward/cbor-gen-for v0.0.0-20200723175505-5892b522820a // indirect
	github.com/ipfs/go-ds-badger2 v0.1.1-0.20200708190120-187fc06f714e // indirect
	github.com/ipfs/go-hamt-ipld v0.1.1 // indirect
	github.com/lib/pq v1.7.0 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/raulk/clock v1.1.0 // indirect
	github.com/renproject/id v0.4.2
	github.com/renproject/pack v0.2.3
	github.com/renproject/surge v1.2.5
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tendermint/tendermint v0.33.8
	github.com/terra-project/core v0.4.0-rc.4
	github.com/xorcare/golden v0.6.1-0.20191112154924-b87f686d7542 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)

replace github.com/cosmos/ledger-cosmos-go => github.com/terra-project/ledger-terra-go v0.11.1-terra

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/CosmWasm/go-cosmwasm => github.com/terra-project/go-cosmwasm v0.10.1-terra
