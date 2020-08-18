package terra

import "github.com/renproject/multichain/chain/cosmos"

type (
	AddressDecoder       = cosmos.AddressDecoder
	AddressEncoder       = cosmos.AddressEncoder
	AddressEncodeDecoder = cosmos.AddressEncodeDecoder
)

var (
	NewAddressDecoder      = cosmos.NewAddressDecoder
	NewAddressEncoder      = cosmos.NewAddressEncoder
	NewAddressEnodeDecoder = cosmos.NewAddressEncodeDecoder
)
