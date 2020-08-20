package celo

import "github.com/renproject/multichain/chain/ethereum"

// An Address on the Celo chain is functionally identical to an address on the
// Ethereum chain.
type Address = ethereum.Address

// An AddressEncoder on the Celo chain is functionally identical to an encoder
// on the Ethereum chain.
type AddressEncoder = ethereum.AddressEncoder

// An AddressDecoder on the Celo chain is functionally identical to a decoder on
// the Ethereum chain.
type AddressDecoder = ethereum.AddressDecoder

// An AddressEncoderDecoder on the Celo chain is functionally identical to a
// encoder/decoder on the Ethereum chain.
type AddressEncoderDecoder = ethereum.AddressEncoderDecoder

var (
	NewAddressEncoderDecoder = ethereum.NewAddressEncoderDecoder
)
