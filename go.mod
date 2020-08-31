module github.com/renproject/multichain

go 1.14

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/centrifuge/go-substrate-rpc-client v1.1.0
	github.com/codahale/blake2 v0.0.0-20150924215134-8d10d0420cbf
	github.com/ethereum/go-ethereum v1.9.20
	github.com/filecoin-project/go-address v0.0.3
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20200822201400-474f4fdccc52
	github.com/filecoin-project/lotus v0.5.6
	github.com/filecoin-project/specs-actors v0.9.3
	github.com/ipfs/go-cid v0.0.7
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-varint v0.0.6
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/renproject/id v0.4.2
	github.com/renproject/pack v0.2.3
	github.com/renproject/surge v1.2.6
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)

replace github.com/filecoin-project/filecoin-ffi => ./chain/filecoin/filecoin-ffi
