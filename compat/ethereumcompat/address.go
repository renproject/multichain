package ethereumcompat

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/renproject/pack"
	"github.com/renproject/surge"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type AddressDecoder interface {
	DecodeAddress(pack.String) (Address, error)
}

type addressDecoder struct{}

func NewAddressDecoder() AddressDecoder {
	return addressDecoder{}
}

func (addressDecoder) DecodeAddress(encoded pack.String) (Address, error) {
	return NewAddressFromHex(encoded.String())
}

// An Address represents a public address on the Ethereum blockchain. It can be
// the address of an external account, or the address of a smart contract.
type Address ethcommon.Address

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
	ethaddr := ethcommon.Address{}
	copy(ethaddr[:], ethaddrData)
	return Address(ethaddr), nil
}

// SizeHint returns the number of bytes needed to represent this address in
// binary.
func (Address) SizeHint() int {
	return ethcommon.AddressLength
}

// Marshal the address to binary.
func (addr Address) Marshal(buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < ethcommon.AddressLength || rem < ethcommon.AddressLength {
		return buf, rem, surge.ErrUnexpectedEndOfBuffer
	}
	copy(buf, addr[:])
	return buf[ethcommon.AddressLength:], rem - ethcommon.AddressLength, nil
}

// Unmarshal the address from binary.
func (addr *Address) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < ethcommon.AddressLength || rem < ethcommon.AddressLength {
		return buf, rem, surge.ErrUnexpectedEndOfBuffer
	}
	copy(addr[:], buf[:ethcommon.AddressLength])
	return buf[ethcommon.AddressLength:], rem - ethcommon.AddressLength, nil
}

// MarshalJSON implements JSON marshaling by encoding the address as a hex
// string.
func (addr Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(ethcommon.Address(addr).Hex())
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
