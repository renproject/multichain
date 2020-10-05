package icon

import (
	"encoding/hex"
	"log"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/icon/crypto"
	"github.com/renproject/pack"
)

// An Address is a public address that can be encoded/decoded to/from strings.
// Addresses are usually formatted different between different network
// configurations.
type Address [AddressBytes]byte

const (
	// AddressIDBytes represents an Icon address body size
	AddressIDBytes = 20
	// AddressBytes represents an Icon address byte size
	AddressBytes = AddressIDBytes + 1
)

// AddressEncodeDecoder encapsulates fields that implement the
// address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// AddressEncoder implements the address.Encoder interface
type AddressEncoder struct{}

// AddressDecoder implements the address.Decoder interface
type AddressDecoder struct{}

// NewAddressEncodeDecoder constructs a new EncodeDecoder.
func NewAddressEncodeDecoder() address.EncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(),
		AddressDecoder: NewAddressDecoder(),
	}
}

// NewAddressDecoder constructs a new AddressDecoder.
func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

// NewAddressEncoder constructs a new AddressEncoder.
func NewAddressEncoder() AddressEncoder {
	return AddressEncoder{}
}

// EncodeAddress consumes raw bytes and encodes them to a human-readable
// address format.
func (ae AddressEncoder) EncodeAddress(raw address.RawAddress) (address.Address, error) {
	post := hex.EncodeToString(raw[1:])
	if raw[0] == 1 {
		return address.Address("cx" + post), nil
	}
	return address.Address("hx" + post), nil
}

// DecodeAddress consumes a human-readable representation of an icon
// compatible address and decodes it to its raw bytes representation.
func (de AddressDecoder) DecodeAddress(ai address.Address) (address.RawAddress, error) {
	var isContract = false
	if len(ai) >= 2 {
		switch {
		case ai[0:2] == "cx":
			isContract = true
			ai = ai[2:]
		case ai[0:2] == "hx":
			ai = ai[2:]
		case ai[0:2] == "0x":
			ai = ai[2:]
		}
	}
	a, err := hex.DecodeString(string(ai))
	if err != nil {
		return a, err
	}
	res := make(address.RawAddress, len(a)+1)
	copy(res[1:], a)
	if isContract {
		res[0] = 1
	} else {
		res[0] = 0
	}
	return address.RawAddress(pack.Bytes(res)), nil
}

// NewAccountAddressFromPublicKey generates an address from a public key
func NewAccountAddressFromPublicKey(pubKey *crypto.PublicKey) *Address {
	a := new(Address)
	pk := pubKey.SerializeUncompressed()
	if pk == nil {
		log.Panicln("FAIL invalid public key:", pubKey)
	}
	digest := crypto.SHA3Sum256(pk[1:])
	a.SetTypeAndID(false, digest[len(digest)-20:])
	return a
}

// SetTypeAndID generates an address from a public key
func (a *Address) SetTypeAndID(ic bool, id []byte) error {
	if id == nil {
		return nil
	}
	switch {
	case len(id) < AddressIDBytes:
		copy(a[AddressIDBytes-len(id)+1:], id)
	default:
		copy(a[1:], id)
	}
	if ic {
		a[0] = 1
	} else {
		a[0] = 0
	}
	return nil
}

// String returns the address as a human-readable hex string.
func (a *Address) String() string {
	if a[0] == 1 {
		return "cx" + hex.EncodeToString(a[1:])
	}
	return "hx" + hex.EncodeToString(a[1:])
}

// Bytes returns the address as a slice of 20 bytes.
func (a *Address) Bytes() []byte {
	return (*a)[:]
}
