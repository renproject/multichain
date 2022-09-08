package digibyte

import "github.com/renproject/multichain/chain/utxochain"

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder = utxochain.AddressEncoder

// AddressDecoder encapsulates the chain specific configurations and implements
// the address.Decoder interface
type AddressDecoder = utxochain.AddressDecoder

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder = utxochain.AddressEncodeDecoder
