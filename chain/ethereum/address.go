package ethereum

import (
	"github.com/renproject/multichain/chain/evm"
)

type (
	// AddressEncodeDecoder re-exports evm.AddressEncodeDecoder.
	AddressEncodeDecoder = evm.AddressEncodeDecoder

	// AddressEncoder re-exports evm.AddressEncoder.
	AddressEncoder = evm.AddressEncoder

	// AddressDecoder re-exports evm.AddressDecoder.
	AddressDecoder = evm.AddressDecoder

	// Address re-exports evm.Address.
	Address = evm.Address
)

var (
	// NewAddressEncodeDecoder re-exports evm.NewAddressEncodeDecoder.
	NewAddressEncodeDecoder = evm.NewAddressEncodeDecoder

	// NewAddressDecoder re-exports evm.NewAddressDecoder.
	NewAddressDecoder = evm.NewAddressDecoder

	// NewAddressEncoder re-exports evm.NewAddressEncoder.
	NewAddressEncoder = evm.NewAddressEncoder

	// NewAddressFromHex re-exports evm.NewAddressFromHex.
	NewAddressFromHex = evm.NewAddressFromHex
)
