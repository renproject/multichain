package intconv

import (
	"math/big"
)

var BigIntOne = big.NewInt(1)

func BytesForZero() []byte {
	return []byte{0}
}

func BigIntToBytes(i *big.Int) []byte {
	if i == nil || i.Sign() == 0 {
		return BytesForZero()
	} else if i.Sign() > 0 {
		bl := i.BitLen()
		if (bl % 8) == 0 {
			bs := make([]byte, bl/8+1)
			copy(bs[1:], i.Bytes())
			return bs
		}
		return i.Bytes()
	} else {
		var ti, nb big.Int
		ti.Add(i, BigIntOne)
		bl := ti.BitLen()
		nb.SetBit(&nb, (bl+8)/8*8, 1)
		nb.Add(&nb, i)
		return nb.Bytes()
	}
}

func BigIntSetBytes(i *big.Int, bs []byte) *big.Int {
	i.SetBytes(bs)
	if len(bs) > 0 && (bs[0]&0x80) != 0 {
		var base big.Int
		base.SetBit(&base, i.BitLen(), 1)
		i.Sub(i, &base)
	}
	return i
}

func Uint64ToBytes(v uint64) []byte {
	if v == 0 {
		return BytesForZero()
	}
	bs := make([]byte, 9)
	for idx := 8; idx >= 0; idx-- {
		tv := byte(v & 0xff)
		bs[idx] = tv
		v >>= 8
		if v == 0 && (tv&0x80) == 0 {
			return bs[idx:]
		}
	}
	return bs
}

func SizeToBytes(v uint64) []byte {
	if v == 0 {
		return BytesForZero()
	}
	bs := make([]byte, 8)
	for idx := 7; idx >= 0; idx-- {
		bs[idx] = byte(v & 0xff)
		v >>= 8
		if v == 0 {
			return bs[idx:]
		}
	}
	return bs
}

func BytesToUint64(bs []byte) uint64 {
	if len(bs) == 0 {
		return 0
	}
	var v uint64
	if (bs[0] & 0x80) != 0 {
		v = 0xffffffffffffffff
	}
	for _, b := range bs {
		v = (v << 8) | uint64(b)
	}
	return v
}

func BytesToSize(bs []byte) uint64 {
	var v uint64
	for _, b := range bs {
		v = (v << 8) | uint64(b)
	}
	return v
}

func BytesToInt64(bs []byte) int64 {
	if len(bs) == 0 {
		return 0
	}
	var v int64
	if (bs[0] & 0x80) != 0 {
		for _, b := range bs {
			v = (v << 8) | int64(b^0xff)
		}
		return -v - 1
	} else {
		for _, b := range bs {
			v = (v << 8) | int64(b)
		}
		return v
	}
}

func Int64ToBytes(v int64) []byte {
	if v == 0 {
		return BytesForZero()
	}
	bs := make([]byte, 8)

	const mask int64 = -0x80
	var target int64 = 0
	if v < 0 {
		target = mask
	}
	for idx := 7; idx >= 0; idx-- {
		bs[idx] = byte(v & 0xff)
		if (v & mask) == target {
			return bs[idx:]
		}
		v >>= 8
	}
	return bs
}
