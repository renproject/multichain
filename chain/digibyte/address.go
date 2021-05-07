package digibyte

import "github.com/renproject/multichain/chain/bitcoin"

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder = bitcoin.AddressEncoder

// AddressDecoder encapsulates the chain specific configurations and implements
// the address.Decoder interface
type AddressDecoder = bitcoin.AddressDecoder

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder = bitcoin.AddressEncodeDecoder
