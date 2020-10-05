package intconv

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
)

func FormatBigInt(i *big.Int) string {
	return encodeHexNumber(i.Sign() < 0, i.Bytes())
}

func ParseBigInt(i *big.Int, s string) error {
	neg, bs, err := decodeHexNumber(s)
	if err != nil {
		return err
	}
	i.SetBytes(bs)
	if neg {
		i.Neg(i)
	}
	return nil
}

func encodeHexNumber(neg bool, b []byte) string {
	s := hex.EncodeToString(b)
	if len(s) == 0 {
		return "0x0"
	}
	if s[0] == '0' {
		s = s[1:]
	}
	if neg {
		return "-0x" + s
	} else {
		return "0x" + s
	}
}

func decodeHexNumber(s string) (bool, []byte, error) {
	negative := false
	if len(s) > 0 && s[0] == '-' {
		negative = true
		s = s[1:]
	}
	if len(s) > 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	if (len(s) % 2) == 1 {
		s = "0" + s
	}
	bs, err := hex.DecodeString(s)
	return negative, bs, err
}

func ParseInt(s string, bits int) (int64, error) {
	if v64, err := strconv.ParseInt(s, 0, bits); err == nil {
		return v64, nil
	}
	if negative, bs, err := decodeHexNumber(s); err == nil {
		if len(bs)*8 > bits {
			return 0, errors.New("OutOfRange")
		}
		u64 := BytesToSize(bs)
		edge := (uint64(1)) << uint(bits-1)
		if negative {
			if u64 > edge {
				return 0, errors.New("OutOfRange")
			}
			return -int64(u64), nil
		} else {
			if u64 >= edge {
				return 0, errors.New("OutOfRange")
			}
			return int64(u64), nil
		}
	} else {
		return 0, err
	}
}

func ParseUint(s string, bits int) (uint64, error) {
	if v64, err := strconv.ParseUint(s, 0, bits); err == nil {
		return v64, nil
	}
	if negative, bs, err := decodeHexNumber(s); err == nil && !negative {
		if len(bs)*8 > bits {
			return 0, errors.New("OutOfRange")
		}
		return BytesToSize(bs), nil
	} else {
		return 0, errors.New("IllegalFormat")
	}
}

func FormatInt(v int64) string {
	var bs []byte
	if v < 0 {
		bs = SizeToBytes(uint64(-v))
		return encodeHexNumber(true, bs)
	} else {
		bs = SizeToBytes(uint64(v))
		return encodeHexNumber(false, bs)
	}
}

func FormatUint(v uint64) string {
	return encodeHexNumber(false, SizeToBytes(v))
}
