package ethereum

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// AddressEncoder implements the address.Encoder interface.
type AddressEncoder interface {
	EncodeAddress(address.RawAddress) (address.Address, error)
}

type addressEncoder struct{}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder.
func NewAddressEncodeDecoder() address.EncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(),
		AddressDecoder: NewAddressDecoder(),
	}
}

// AddressDecoder implements the address.Decoder interface.
type AddressDecoder interface {
	DecodeAddress(address.Address) (address.RawAddress, error)
}

type addressDecoder struct{}

// NewAddressDecoder constructs a new AddressDecoder.
func NewAddressDecoder() AddressDecoder {
	return addressDecoder{}
}

// NewAddressEncoder constructs a new AddressEncoder.
func NewAddressEncoder() AddressEncoder {
	return addressEncoder{}
}

func (addressDecoder) DecodeAddress(encoded address.Address) (address.RawAddress, error) {
	ethaddr, err := NewAddressFromHex(string(pack.String(encoded)))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(pack.Bytes(ethaddr[:])), nil
}

func (addressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	encodedAddr := common.Bytes2Hex([]byte(rawAddr))
	return address.Address(pack.NewString(encodedAddr)), nil
}

// An Address represents a public address on the Ethereum blockchain. It can be
// the address of an external account, or the address of a smart contract.
type Address common.Address

// NewAddressFromHex returns an Address decoded from a hex
// string.
func NewAddressFromHex(str string) (Address, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	if len(str) != 40 {
		return Address{}, fmt.Errorf("invalid ethaddress %v", str)
	}
	ethaddrData, err := hex.DecodeString(str)
	if err != nil {
		return Address{}, fmt.Errorf("invalid ethaddress %v: %v", str, err)
	}
	ethaddr := common.Address{}
	copy(ethaddr[:], ethaddrData)
	return Address(ethaddr), nil
}

// SizeHint returns the number of bytes needed to represent this address in
// binary.
func (Address) SizeHint() int {
	return common.AddressLength
}

// Marshal the address to binary.
func (addr Address) Marshal(buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < common.AddressLength || rem < common.AddressLength {
		return buf, rem, surge.ErrUnexpectedEndOfBuffer
	}
	copy(buf, addr[:])
	return buf[common.AddressLength:], rem - common.AddressLength, nil
}

// Unmarshal the address from binary.
func (addr *Address) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < common.AddressLength || rem < common.AddressLength {
		return buf, rem, surge.ErrUnexpectedEndOfBuffer
	}
	copy(addr[:], buf[:common.AddressLength])
	return buf[common.AddressLength:], rem - common.AddressLength, nil
}

// MarshalJSON implements JSON marshaling by encoding the address as a hex
// string.
func (addr Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(common.Address(addr).Hex())
}

// UnmarshalJSON implements JSON unmarshaling by expected the data be a hex
// encoded string representation of an address.
func (addr *Address) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ethaddr, err := NewAddressFromHex(str)
	if err != nil {
		return err
	}
	*addr = ethaddr
	return nil
}

// String returns the address as a human-readable hex string.
func (addr Address) String() string {
	return hex.EncodeToString(addr[:])
}

// Bytes returns the address as a slice of 20 bytes.
func (addr Address) Bytes() pack.Bytes {
	return pack.Bytes(addr[:])
}
