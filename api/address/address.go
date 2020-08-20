// Package address defines the Address API. All chains must implement this API,
// so that addresses can be encoded/decoded.
package address

import "github.com/renproject/pack"

// An Address is a human-readable representation of a public identity. It can be
// the address of an external account, contract, or script.
type Address pack.String

// SizeHint returns the number of bytes required to represent the address in
// binary.
func (addr Address) SizeHint() int {
	return pack.String(addr).SizeHint()
}

// Marshal the address to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (addr Address) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.String(addr).Marshal(buf, rem)
}

// Unmarshal the address from binary. You should not call this function
// directly, unless you are implementing unmarshalling for a container type.
func (addr *Address) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.String)(addr).Unmarshal(buf, rem)
}

// RawAddress is an address that has been decoded into its binary form.
type RawAddress pack.Bytes

// SizeHint returns the number of bytes required to represent the address in
// binary.
func (addr RawAddress) SizeHint() int {
	return pack.Bytes(addr).SizeHint()
}

// Marshal the address to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (addr RawAddress) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.Bytes(addr).Marshal(buf, rem)
}

// Unmarshal the address from binary. You should not call this function
// directly, unless you are implementing unmarshalling for a container type.
func (addr *RawAddress) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.Bytes)(addr).Unmarshal(buf, rem)
}

// The Encoder interface is used to convert raw addresses into human-readable
// addresses.
type Encoder interface {
	EncodeAddress(RawAddress) (Address, error)
}

// The Decoder interfaces is used to convert human-readable addresses into raw
// addresses.
type Decoder interface {
	DecodeAddress(Address) (RawAddress, error)
}

// The EncoderDecoder interfaces combines encoding and decoding functionality
// into one interface.
type EncoderDecoder interface {
	Encoder
	Decoder
}
