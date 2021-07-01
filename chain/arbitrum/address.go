package arbitrum

import (
	"github.com/renproject/multichain/chain/ethereum"
)

type (
	// AddressEncodeDecoder re-exports ethereum.AddressEncodeDecoder.
	AddressEncodeDecoder = ethereum.AddressEncodeDecoder

	// AddressEncoder re-exports ethereum.AddressEncoder.
	AddressEncoder = ethereum.AddressEncoder

	// AddressDecoder re-exports ethereum.AddressDecoder.
	AddressDecoder = ethereum.AddressDecoder

	// Address re-exports ethereum.Address.
	Address = ethereum.Address
)

var (
	// NewAddressEncodeDecoder re-exports ethereum.NewAddressEncodeDecoder.
	NewAddressEncodeDecoder = ethereum.NewAddressEncodeDecoder

	// NewAddressDecoder re-exports ethereum.NewAddressDecoder.
	NewAddressDecoder = ethereum.NewAddressDecoder

	// NewAddressEncoder re-exports ethereum.NewAddressEncoder.
	NewAddressEncoder = ethereum.NewAddressEncoder

	// NewAddressFromHex re-exports ethereum.NewAddressFromHex.
	NewAddressFromHex = ethereum.NewAddressFromHex
)
