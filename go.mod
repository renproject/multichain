module github.com/renproject/multichain

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.44.0
	github.com/dchest/blake2b v1.0.0
	github.com/ethereum/go-ethereum v1.10.6
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/go-state-types v0.1.1-0.20210506134452-99b279731c48
	github.com/filecoin-project/lotus v1.10.0
	github.com/ipfs/go-cid v0.0.7
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-varint v0.0.6
	github.com/near/borsh-go v0.3.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/renproject/id v0.4.2
	github.com/renproject/pack v0.2.5
	github.com/renproject/solana-ffi v0.1.2
	github.com/renproject/surge v1.2.6
	github.com/tendermint/tendermint v0.34.12
	github.com/terra-money/core v0.5.5
	github.com/tyler-smith/go-bip39 v1.1.0
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/filecoin-project/filecoin-ffi => ./chain/filecoin/filecoin-ffi

replace github.com/renproject/solana-ffi => ./chain/solana/solana-ffi

replace github.com/cosmos/ledger-cosmos-go => github.com/terra-money/ledger-terra-go v0.11.2

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.34.12

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
