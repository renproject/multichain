package crypto

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/haltingstate/secp256k1-go"
)

const (
	// PrivateKeyLen is the byte length of a private key
	PrivateKeyLen = 32
)

// PrivateKey is a type representing a private key.
// TODO private key always includes public key? or create KeyPair struct
// for both private key and public key
type PrivateKey struct {
	bytes []byte // 32-byte
}

// String returns the string representation.
func (key *PrivateKey) String() string {
	return hex.EncodeToString(key.bytes)
}

// PublicKey generates a public key paired with itself.
func (key *PrivateKey) PublicKey() *PublicKey {
	pkBytes := secp256k1.PubkeyFromSeckey(key.bytes)
	pk, err := ParsePublicKey(pkBytes)
	if err != nil {
		panic(err)
	}
	return pk
}

// Bytes returns bytes form of private key.
func (key *PrivateKey) Bytes() []byte {
	kb := make([]byte, PrivateKeyLen)
	copy(kb, key.bytes)
	return kb
}

// TODO add 'func ToECDSA() ecdsa.PrivateKey' if needed

const (
	// PublicKeyLenCompressed is the byte length of a compressed public key
	PublicKeyLenCompressed = 33
	// PublicKeyLenUncompressed is the byte length of an uncompressed public key
	PublicKeyLenUncompressed = 65

	publicKeyCompressed   byte = 0x2 // y_bit + x coord
	publicKeyUncompressed byte = 0x4 // x coord + y coord
)

// PublicKey is a type representing a public key, which can be serialized to
// or deserialized from compressed or uncompressed formats.
type PublicKey struct {
	bytes []byte // 33-byte compressed format to use halting state library efficiently
}

// ParsePublicKey parses the public key into a PublicKey instance. It supports
// uncompressed and compressed formats.
// NOTE: For the efficiency, it may use the slice directly. So don't change any
// internal value of the public key
func ParsePublicKey(pubKey []byte) (*PublicKey, error) {
	switch len(pubKey) {
	case 0:
		return nil, errors.New("public key bytes are empty")
	case PublicKeyLenCompressed:
		return &PublicKey{pubKey}, nil
	case PublicKeyLenUncompressed:
		return &PublicKey{uncompToCompPublicKey(pubKey)}, nil
	default:
		return nil, errors.New("wrong format")
	}
}

// uncompToCompPublicKey changes the uncompressed formatted public key to
// the compressed formatted. It assumes the uncompressed key is valid.
func uncompToCompPublicKey(uncomp []byte) (comp []byte) {
	comp = make([]byte, PublicKeyLenCompressed)
	// skip to check the validity of uncompressed key
	format := publicKeyCompressed
	if uncomp[64]&0x1 == 0x1 {
		format |= 0x1
	}
	comp[0] = format
	copy(comp[1:], uncomp[1:33])
	return
}

// SerializeCompressed serializes the public key in a 33-byte compressed format.
// For the efficiency, it returns the slice internally used, so don't change
// any internal value in the returned slice.
func (key *PublicKey) SerializeCompressed() []byte {
	return key.bytes
}

// SerializeUncompressed serializes the public key in a 65-byte uncompressed format.
func (key *PublicKey) SerializeUncompressed() []byte {
	return secp256k1.UncompressPubkey(key.bytes)
}

// Equal returns true if the given public key is same as this instance
// semantically
func (key *PublicKey) Equal(key2 *PublicKey) bool {
	return bytes.Equal(key.bytes, key2.bytes)
}

// String returns the string representation.
func (key *PublicKey) String() string {
	return "0x" + hex.EncodeToString(key.bytes)
}

// GenerateKeyPair generates a private and public key pair.
func GenerateKeyPair() (privKey *PrivateKey, pubKey *PublicKey) {
	pub, priv := secp256k1.GenerateKeyPair()
	privKey = &PrivateKey{priv}
	pubKey, _ = ParsePublicKey(pub)
	return
}

// ParsePublicKey parse private key and return private key object.
func ParsePrivateKey(b []byte) (*PrivateKey, error) {
	if len(b) != PrivateKeyLen {
		return nil, errors.New("InvalidKeyLength")
	}
	b2 := make([]byte, len(b))
	copy(b2, b)
	return &PrivateKey{b2}, nil
}
