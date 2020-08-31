module github.com/renproject/multichain

go 1.14

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/centrifuge/go-substrate-rpc-client v1.1.0
	github.com/codahale/blake2 v0.0.0-20150924215134-8d10d0420cbf
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/drand/drand v1.0.3-0.20200714175734-29705eaf09d4 // indirect
	github.com/elastic/go-sysinfo v1.4.0 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/ethereum/go-ethereum v1.9.19
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-amt-ipld v0.0.0-20191205011053-79efc22d6cdc // indirect
	github.com/filecoin-project/go-amt-ipld/v2 v2.1.0 // indirect
	github.com/filecoin-project/go-data-transfer v0.5.0 // indirect
	github.com/filecoin-project/go-fil-markets v0.3.2-0.20200702145639-4034a18364e4
	github.com/filecoin-project/lotus v0.4.1
	github.com/filecoin-project/sector-storage v0.0.0-20200723200950-ed2e57dde6df // indirect
	github.com/filecoin-project/specs-actors v0.6.2-0.20200702170846-2cd72643a5cf
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-ds-badger2 v0.1.1-0.20200708190120-187fc06f714e // indirect
	github.com/ipfs/go-hamt-ipld v0.1.1 // indirect
	github.com/ipld/go-ipld-prime v0.0.3 // indirect
	github.com/lib/pq v1.7.0 // indirect
	github.com/libp2p/go-libp2p-core v0.6.1 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-multiaddr v0.3.1 // indirect
	github.com/multiformats/go-varint v0.0.6
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/raulk/clock v1.1.0 // indirect
	github.com/renproject/id v0.4.2
	github.com/renproject/pack v0.2.3
	github.com/renproject/surge v1.2.5
	github.com/tendermint/tendermint v0.33.8
	github.com/terra-project/core v0.3.7
	github.com/whyrusleeping/cbor-gen v0.0.0-20200715143311-227fab5a2377
	github.com/xorcare/golden v0.6.1-0.20191112154924-b87f686d7542 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.0.0-20200824131525-c12d262b63d8 // indirect
	golang.org/x/tools v0.0.0-20200825202427-b303f430e36d // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	howett.net/plist v0.0.0-20200419221736-3b63eb3a43b5 // indirect
	modernc.org/golex v1.0.1 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./chain/filecoin/filecoin-ffi
