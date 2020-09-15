package terra

import "github.com/renproject/multichain/chain/cosmos"

type (
	// Address re-exports cosmos-compatible address
	Address = cosmos.Address

	// AddressDecoder re-exports cosmos.AddressDecoder
	AddressDecoder = cosmos.AddressDecoder

	// AddressEncoder re-exports cosmos.AddressEncoder
	AddressEncoder = cosmos.AddressEncoder

	// AddressEncodeDecoder re-exports cosmos.AddressEncodeDecoder
	AddressEncodeDecoder = cosmos.AddressEncodeDecoder
)

var (
	// NewAddressDecoder re-exports cosmos.NewAddressDecoder
	NewAddressDecoder = cosmos.NewAddressDecoder

	// NewAddressEncoder re-exports cosmos.NewAddressEncoder
	NewAddressEncoder = cosmos.NewAddressEncoder

	// NewAddressEncodeDecoder re-exports cosmos.NewAddressEnodeDecoder
	NewAddressEncodeDecoder = cosmos.NewAddressEncodeDecoder
)
