package zcash

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// An Address represents a Zcash address.
type Address interface {
	btcutil.Address
	BitcoinAddress() btcutil.Address
}

// AddressPubKeyHash represents an address for P2PKH transactions for Zcash that
// is compatible with the Bitcoin Compat API.
type AddressPubKeyHash struct {
	*btcutil.AddressPubKeyHash
	params *Params
}

// NewAddressPubKeyHash returns a new AddressPubKeyHash that is compatible with
// the Bitcoin Compat API.
func NewAddressPubKeyHash(pkh []byte, params *Params) (AddressPubKeyHash, error) {
	addr, err := btcutil.NewAddressPubKeyHash(pkh, params.Params)
	return AddressPubKeyHash{AddressPubKeyHash: addr, params: params}, err
}

// String returns the string encoding of the transaction output destination.
//
// Please note that String differs subtly from EncodeAddress: String will return
// the value as a string without any conversion, while EncodeAddress may convert
// destination types (for example, converting pubkeys to P2PKH addresses) before
// encoding as a payment address string.
func (addr AddressPubKeyHash) String() string {
	return addr.EncodeAddress()
}

// EncodeAddress returns the string encoding of the payment address associated
// with the Address value. See the comment on String for how this method differs
// from String.
func (addr AddressPubKeyHash) EncodeAddress() string {
	hash := *addr.AddressPubKeyHash.Hash160()
	return encodeAddress(hash[:], addr.params.P2PKHPrefix)
}

// ScriptAddress returns the raw bytes of the address to be used when inserting
// the address into a txout's script.
func (addr AddressPubKeyHash) ScriptAddress() []byte {
	return addr.AddressPubKeyHash.ScriptAddress()
}

// IsForNet returns whether or not the address is associated with the passed
// bitcoin network.
func (addr AddressPubKeyHash) IsForNet(params *chaincfg.Params) bool {
	return addr.AddressPubKeyHash.IsForNet(params)
}

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressPubKeyHash) BitcoinAddress() btcutil.Address {
	return addr.AddressPubKeyHash
}

// AddressScriptHash represents an address for P2SH transactions for Zcash that
// is compatible with the Bitcoin Compat API.
type AddressScriptHash struct {
	*btcutil.AddressScriptHash
	params *Params
}

// NewAddressScriptHash returns a new AddressScriptHash that is compatible with
// the Bitcoin Compat API.
func NewAddressScriptHash(script []byte, params *Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHash(script, params.Params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// NewAddressScriptHashFromHash returns a new AddressScriptHash that is compatible with
// the Bitcoin Compat API.
func NewAddressScriptHashFromHash(scriptHash []byte, params *Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHashFromHash(scriptHash, params.Params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// String returns the string encoding of the transaction output destination.
//
// Please note that String differs subtly from EncodeAddress: String will return
// the value as a string without any conversion, while EncodeAddress may convert
// destination types (for example, converting pubkeys to P2PKH addresses) before
// encoding as a payment address string.
func (addr AddressScriptHash) String() string {
	return addr.EncodeAddress()
}

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressScriptHash) BitcoinAddress() btcutil.Address {
	return addr.AddressScriptHash
}

// EncodeAddress returns the string encoding of the payment address associated
// with the Address value. See the comment on String for how this method differs
// from String.
func (addr AddressScriptHash) EncodeAddress() string {
	hash := *addr.AddressScriptHash.Hash160()
	return encodeAddress(hash[:], addr.params.P2SHPrefix)
}

// ScriptAddress returns the raw bytes of the address to be used when inserting
// the address into a txout's script.
func (addr AddressScriptHash) ScriptAddress() []byte {
	return addr.AddressScriptHash.ScriptAddress()
}

// IsForNet returns whether or not the address is associated with the passed
// bitcoin network.
func (addr AddressScriptHash) IsForNet(params *chaincfg.Params) bool {
	return addr.AddressScriptHash.IsForNet(params)
}

func DecodeAddress(addr string) (Address, error) {
	var decoded = base58.Decode(addr)
	if len(decoded) != 26 && len(decoded) != 25 {
		return nil, base58.ErrInvalidFormat
	}

	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, base58.ErrChecksum
	}

	if len(decoded)-6 != ripemd160.Size && len(decoded)-5 != ripemd160.Size {
		return nil, errors.New("incorrect payload len")
	}

	var addrType uint8
	var params *Params
	var err error
	var hash [20]byte
	if len(decoded) == 26 {
		addrType, params, err = parsePrefix(decoded[:2])
		copy(hash[:], decoded[2:22])
	} else {
		addrType, params, err = parsePrefix(decoded[:1])
		copy(hash[:], decoded[1:21])
	}
	if err != nil {
		return nil, err
	}

	switch addrType {
	case 0: // P2PKH
		return NewAddressPubKeyHash(hash[:], params)
	case 1: // P2SH
		return NewAddressScriptHashFromHash(hash[:], params)
	}

	return nil, errors.New("unknown address")
}

func encodeAddress(hash, prefix []byte) string {
	var (
		body  = append(prefix, hash...)
		chk   = checksum(body)
		cksum [4]byte
	)
	copy(cksum[:], chk[:4])
	return base58.Encode(append(body, cksum[:]...))
}

func checksum(input []byte) (cksum [4]byte) {
	var (
		h  = sha256.Sum256(input)
		h2 = sha256.Sum256(h[:])
	)
	copy(cksum[:], h2[:4])
	return
}

func parsePrefix(prefix []byte) (uint8, *Params, error) {
	if bytes.Equal(prefix, MainNetParams.P2PKHPrefix) {
		return 0, &MainNetParams, nil
	}
	if bytes.Equal(prefix, MainNetParams.P2SHPrefix) {
		return 1, &MainNetParams, nil
	}
	if bytes.Equal(prefix, TestNet3Params.P2PKHPrefix) {
		return 0, &TestNet3Params, nil
	}
	if bytes.Equal(prefix, TestNet3Params.P2SHPrefix) {
		return 1, &TestNet3Params, nil
	}
	if bytes.Equal(prefix, RegressionNetParams.P2PKHPrefix) {
		return 0, &RegressionNetParams, nil
	}
	if bytes.Equal(prefix, RegressionNetParams.P2SHPrefix) {
		return 1, &RegressionNetParams, nil
	}
	return 0, nil, btcutil.ErrUnknownAddressType
}
