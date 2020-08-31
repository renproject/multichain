// Package contract defines the Contract API. All chains that have "smart
// contracts" should implement this API. UTXO-based chains that support
// scripting must not attempt to implementing scripting using this API.
package contract

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

// CallData is used to specify a function and its parameters when invoking
// business logic on a contract.
type CallData pack.Bytes

// SizeHint returns the number of bytes required to represent the calldata in
// binary.
func (data CallData) SizeHint() int {
	return pack.Bytes(data).SizeHint()
}

// Marshal the address to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (data CallData) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.Bytes(data).Marshal(buf, rem)
}

// Unmarshal the address from binary. You should not call this function
// directly, unless you are implementing unmarshalling for a container type.
func (data *CallData) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.Bytes)(data).Unmarshal(buf, rem)
}

// The Caller interface defines the functionality required to call readonly
// functions on a contract. Calling functions that mutate contract state should
// be done using the Account API.
type Caller interface {
	// CallContract at the specified address, using the specified calldata as
	// input (this encodes the function and its parameters). The function output
	// is returned as raw uninterpreted bytes. It is up to the application to
	// interpret these bytes in a meaningful way.  If the call cannot be
	// completed before the context is done, or the call is invalid, then an
	// error should be returned.
	CallContract(context.Context, address.Address, CallData) (pack.Bytes, error)
}
