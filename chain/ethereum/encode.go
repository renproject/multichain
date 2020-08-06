package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/renproject/pack"
)

// Encode a pack-encoded value into an RLP-encoded value.
func Encode(input pack.Value) (abi.Arguments, []interface{}, error) {
	switch input := input.(type) {
	case pack.Struct:
		return EncodeTuple(input)
	case pack.Typed:
		return EncodeTuple(pack.Struct(input))
	default:
		return EncodeUnit(input)
	}
}

func EncodeTuple(input pack.Struct) (abi.Arguments, []interface{}, error) {
	ethargs := make(abi.Arguments, 0, len(input))
	ethvals := make([]interface{}, 0, len(input))

	for _, field := range input {
		fieldargs, fieldvals, err := Encode(field.Value)
		if err != nil {
			return nil, nil, err
		}
		ethargs = append(ethargs, fieldargs...)
		ethvals = append(ethvals, fieldvals...)
	}

	return ethargs, ethvals, nil
}

func EncodeUnit(input pack.Value) (abi.Arguments, []interface{}, error) {
	var ethtype abi.Type
	var ethval interface{}
	var err error

	switch input := input.(type) {
	case pack.U8:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Uint8()
	case pack.U16:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Uint16()
	case pack.U32:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Uint32()
	case pack.U64:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Uint64()
	case pack.U128:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Int()
	case pack.U256:
		ethtype, err = abi.NewType("uint256", "uint256", nil)
		ethval = input.Int()
	case pack.String:
		ethtype, err = abi.NewType("string", "string", nil)
		ethval = string(input)
	case pack.Bytes:
		ethtype, err = abi.NewType("bytes", "bytes", nil)
		ethval = []byte(input)
	case pack.Bytes32:
		ethtype, err = abi.NewType("bytes32", "bytes32", nil)
		ethval = [32]byte(input)
	default:
		return nil, nil, fmt.Errorf("bad type: %v", err)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("bad type: %v", err)
	}
	return []abi.Argument{{Type: ethtype}}, []interface{}{ethval}, nil
}
