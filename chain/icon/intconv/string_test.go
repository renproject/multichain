package intconv

import "testing"

func TestParseUint(t *testing.T) {
	type args struct {
		s    string
		bits int
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"T1", args{"0x0", 16}, 0, false},
		{"T2", args{"0xffff", 16}, 0xffff, false},
		{"T3", args{"0xffffffffffffffff", 64}, 0xffffffffffffffff, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint(tt.args.s, tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	type args struct {
		s    string
		bits int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"T1", args{"0x0", 16}, 0, false},
		{"T2", args{"0x7fff", 16}, 0x7fff, false},
		{"T3", args{"-0x8000", 16}, -0x8000, false},
		{"T4", args{"0xffff", 16}, 0, true},
		{"T5", args{"0x0ffff", 16}, 0, true},
		{"T6", args{"-0x8000000000000000", 64}, -0x8000000000000000, false},
		{"T7", args{"-0x10000000000000000", 64}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt(tt.args.s, tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatInt(t *testing.T) {
	type args struct {
		v int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"T0", args{0x00}, "0x0"},
		{"T1", args{-0x1}, "-0x1"},
		{"T2", args{-0x80}, "-0x80"},
		{"T3", args{0x80}, "0x80"},
		{"T4", args{-0xff}, "-0xff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatInt(tt.args.v); got != tt.want {
				t.Errorf("FormatInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
