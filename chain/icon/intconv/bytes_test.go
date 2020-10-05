package intconv

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBigIntToBytes(t *testing.T) {
	type args struct {
		i *big.Int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"T1", args{big.NewInt(-0x1)}, []byte{0xff}},
		{"T2", args{big.NewInt(-0x7f)}, []byte{0x81}},
		{"T3", args{big.NewInt(0x80)}, []byte{0x00, 0x80}},
		{"T4", args{big.NewInt(-0x80)}, []byte{0x80}},
		{"T5", args{big.NewInt(0)}, []byte{0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BigIntToBytes(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64ToBytes(t *testing.T) {
	type args struct {
		v int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"T1", args{-1}, []byte{0xff}},
		{"T2", args{-0x7f}, []byte{0x81}},
		{"T3", args{0x80}, []byte{0x00, 0x80}},
		{"T4", args{-0x80}, []byte{0x80}},
		{"T5", args{0}, []byte{0x00}},
		{"T6", args{-0x8000000000000000}, []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"T7", args{0x7fffffffffffffff}, []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64ToBytes(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToInt64(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"T1", args{[]byte{}}, 0},
		{"T2", args{[]byte{0x80}}, -0x80},
		{"T3", args{[]byte{0x00, 0x80}}, 0x80},
		{"T4", args{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, -0x8000000000000000},
		{"T5", args{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}, -0x7fffffffffffffff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToInt64(tt.args.bs); got != tt.want {
				t.Errorf("BytesToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64ToBytes(t *testing.T) {
	type args struct {
		v uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"T1", args{0x00}, []byte{0x00}},
		{"T2", args{0x80}, []byte{0x00, 0x80}},
		{"T3", args{0x80123456789abcde}, []byte{0x00, 0x80, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde}},
		{"T4", args{0x7fffffffffffffff}, []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64ToBytes(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToUint64(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"T1", args{[]byte{}}, 0},
		{"T2", args{[]byte{0x80}}, 0xffffffffffffff80},
		{"T3", args{[]byte{0x00, 0x80}}, 0x80},
		{"T4", args{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, 0x8000000000000000},
		{"T5", args{[]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}, 0x7fffffffffffffff},
		{"T6", args{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}, 0xffffffffffffffff},
		{"T7", args{[]byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}, 0xffffffffffffffff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToUint64(tt.args.bs); got != tt.want {
				t.Errorf("BytesToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
