package solana

import (
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"github.com/renproject/solana-ffi/cgo"
)

// UniquePubkey creates an atomically incrementing pubkey used for tests and
// benchmarking purposes.
func UniquePubkey() address.Address {
	pubkey := cgo.UniquePubkey()
	return address.Address(pubkey)
}

// ProgramDerivedAddress derives an address for an account that only the given
// program has the authority to sign. The address is of the same form as a
// Solana pubkey, except they are ensured to not be on the es25519 curve and
// thus have no associated private key. This address is deterministic, based
// upon the program and the seeds slice.
func ProgramDerivedAddress(seeds pack.Bytes, program address.Address) address.Address {
	programDerivedAddressEncoded := cgo.ProgramDerivedAddress(seeds, uint32(len(seeds)), string(program))
	return address.Address(programDerivedAddressEncoded)
}
