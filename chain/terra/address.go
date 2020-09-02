package terra

import "github.com/renproject/multichain/chain/cosmos"

type (
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

	// NewAddressEnodeDecoder re-exports cosmos.NewAddressEnodeDecoder
	NewAddressEnodeDecoder = cosmos.NewAddressEncodeDecoder
)
