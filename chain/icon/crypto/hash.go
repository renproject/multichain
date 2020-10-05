package crypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

// SHA3Sum256 returns the SHA3-256 digest of the data
func SHA3Sum256(m []byte) []byte {
	d := sha3.Sum256(m)
	return d[:]
}

// SHASum256 returns the SHA256 digest of the data
func SHASum256(m []byte) []byte {
	d := sha256.Sum256(m)
	return d[:]
}
