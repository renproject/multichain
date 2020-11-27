package solana

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/chain/solana/solana-ffi/cgo"
	"github.com/renproject/pack"
)

func Hello(name string) string {
	return cgo.Hello(name)
}

func UniquePubkey() pack.Bytes32 {
	pubkeyEncoded := cgo.UniquePubkey()
	pubkeyDecoded := base58.Decode(pubkeyEncoded)
	pubkey32 := pack.Bytes32{}
	copy(pubkey32[:], pubkeyDecoded)
	return pubkey32
}

func ProgramDerivedAddress(seeds []byte, program string) pack.Bytes32 {
	programDerivedAddressEncoded := cgo.ProgramDerivedAddress(seeds, uint32(len(seeds)), program)
	programDerivedAddressDecoded := base58.Decode(programDerivedAddressEncoded)
	pubkey32 := pack.Bytes32{}
	copy(pubkey32[:], programDerivedAddressDecoded)
	return pubkey32
}
