package cosmoscompat

import (
	"fmt"

	"github.com/renproject/pack"
)

// The AddressDecoder defines an interface for decoding string representations
// of Bitcoin address into the Address interface.
type AddressDecoder interface {
	DecodeAddress(pack.String) (Address, error)
}

// AddressDecoderCallbacks implements the AddressDecoder interface by allowing
// users to define closures for all required methods.
type AddressDecoderCallbacks struct {
	DecodeAddressCallback func(pack.String) (Address, error)
}

// DecodeAddress delegates the method to an inner callback. If not callback
// exists, then an error is returned.
func (cbs AddressDecoderCallbacks) DecodeAddress(encoded pack.String) (Address, error) {
	if cbs.DecodeAddressCallback != nil {
		return cbs.DecodeAddressCallback(encoded)
	}
	return nil, fmt.Errorf("not implemented")
}
