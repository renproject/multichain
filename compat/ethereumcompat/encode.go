package ethereumcompat

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/renproject/pack"
)

// A Payload is an Ethereum encoded function call. It includes an ABI, the
// function being called from the ABI, and the data being passed to the
// function.
type Payload struct {
	ABI  pack.Bytes `json:"abi"`
	Fn   pack.Bytes `json:"fn"`
	Data pack.Bytes `json:"data"`
}

// Encode values into an Ethereum ABI compatible byte slice.
func Encode(vals ...interface{}) []byte {
	ethargs := make(abi.Arguments, 0, len(vals))
	ethvals := make([]interface{}, 0, len(vals))

	for _, val := range vals {
		var ethval interface{}
		var ty abi.Type
		var err error

		switch val := val.(type) {
		case pack.Bytes:
			println("BYTES")
			ethval = []byte(val)
			ty, err = abi.NewType("bytes", "", nil)
		case pack.Bytes32:
			ethval = [32]byte(val)
			ty, err = abi.NewType("bytes32", "", nil)
		case pack.U8:
			ethval = big.NewInt(0).SetUint64(uint64(val.Uint8()))
			ty, err = abi.NewType("uint256", "", nil)
		case pack.U16:
			ethval = big.NewInt(0).SetUint64(uint64(val.Uint16()))
			ty, err = abi.NewType("uint256", "", nil)
		case pack.U32:
			ethval = big.NewInt(0).SetUint64(uint64(val.Uint32()))
			ty, err = abi.NewType("uint256", "", nil)
		case pack.U64:
			ethval = big.NewInt(0).SetUint64(uint64(val.Uint64()))
			ty, err = abi.NewType("uint256", "", nil)
		case pack.U128:
			ethval = val.Int()
			ty, err = abi.NewType("uint256", "", nil)
		case pack.U256:
			ethval = val.Int()
			ty, err = abi.NewType("uint256", "", nil)
		case Address:
			ethval = val
			ty, err = abi.NewType("address", "", nil)
		default:
			panic(fmt.Errorf("non-exhaustive pattern: %T", val))
		}

		if err != nil {
			panic(fmt.Errorf("error encoding: %v", err))
		}
		ethargs = append(ethargs, abi.Argument{
			Type: ty,
		})
		ethvals = append(ethvals, ethval)
	}

	packed, err := ethargs.Pack(ethvals...)
	if err != nil {
		panic(fmt.Errorf("error packing: %v", err))
	}
	return packed
}
