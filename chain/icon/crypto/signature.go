package crypto

import (
	"encoding/hex"
	"errors"
	"sync"

	"github.com/haltingstate/secp256k1-go"
)

const (
	// SignatureLenRawWithV is the bytes length of signature including V value
	SignatureLenRawWithV = 65
	// SignatureLenRaw is the bytes length of signature not including V value
	SignatureLenRaw = 64
	invalidV        = 0xff
	// HashLen is the bytes length of hash for signature
	HashLen = 32
)

// Signature is a type representing an ECDSA signature with or without V.
type Signature struct {
	bytes []byte // 65 bytes of [R|S|V]
	hasV  bool
}

var globalLock sync.Mutex

// NewSignature calculates an ECDSA signature including V, which is 0 or 1.
func NewSignature(hash []byte, privKey *PrivateKey) (*Signature, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if len(hash) == 0 || len(hash) > HashLen || privKey == nil {
		return nil, errors.New("Invalid arguments")
	}
	return ParseSignature(secp256k1.Sign(hash, privKey.bytes))
}

// ParseSignature parses a signature from the raw byte array of 64([R|S]) or
// 65([R|S|V]) bytes long. If a source signature is formatted as [V|R|S],
// call ParseSignatureVRS instead.
// NOTE: For the efficiency, it may use the slice directly. So don't change any
// internal value of the signature.
func ParseSignature(sig []byte) (*Signature, error) {
	var s Signature
	switch len(sig) {
	case 0:
		return nil, errors.New("sigature bytes are empty")
	case SignatureLenRawWithV:
		s.bytes = sig
		s.hasV = true
	case SignatureLenRaw:
		s.bytes = append(s.bytes, sig...)
		s.bytes = append(s.bytes, 0x00) // no meaning
		s.hasV = false
	default:
		return nil, errors.New("wrong raw signature format")
	}
	return &s, nil
}

// ParseSignatureVRS parses a signature from the [V|R|S] formatted signature.
// If the format of a source signature is different,
// call ParseSignature instead.
func ParseSignatureVRS(sig []byte) (*Signature, error) {
	if len(sig) != 65 {
		return nil, errors.New("wrong raw signature format")
	}

	var s Signature
	s.bytes = append(s.bytes, sig[1:33]...)
	s.bytes = append(s.bytes, sig[33:]...)
	s.bytes[64] = sig[0]
	return &s, nil
}

// HasV returns whether the signature has V value.
func (sig *Signature) HasV() bool {
	return sig.hasV
}

// SerializeRS returns the 64-byte data formatted as [R|S] from the signature.
// For the efficiency, it returns the slice internally used, so don't change
// any internal value in the returned slice.
func (sig *Signature) SerializeRS() ([]byte, error) {
	if len(sig.bytes) < 64 {
		return nil, errors.New("not a valid signature")
	}
	return sig.bytes[:64], nil
}

// SerializeVRS returns the 65-byte data formatted as [V|R|S] from the signature.
// Make sure that it has a valid V value. If it doesn't have V value, then it
// will throw error.
// For the efficiency, it returns the slice internally used, so don't change
// any internal value in the returned slice.
func (sig *Signature) SerializeVRS() ([]byte, error) {
	if !sig.HasV() {
		return nil, errors.New("no V value")
	}

	s := make([]byte, SignatureLenRawWithV)
	s[0] = sig.bytes[64]
	copy(s[1:33], sig.bytes[:32])
	copy(s[33:], sig.bytes[32:64])
	return s, nil
}

// SerializeRSV returns the 65-byte data formatted as [R|S|V] from the signature.
// Make sure that it has a valid V value. If it doesn't have V value, then it
// will throw error.
// For the efficiency, it returns the slice internally used, so don't change
// any internal value in the returned slice.
func (sig *Signature) SerializeRSV() ([]byte, error) {
	if !sig.HasV() {
		return nil, errors.New("no V value")
	}

	return sig.bytes, nil
}

// RecoverPublicKey recovers a public key from the hash of message and its signature.
func (sig *Signature) RecoverPublicKey(hash []byte) (*PublicKey, error) {
	if !sig.HasV() {
		return nil, errors.New("signature has no V value")
	}
	if len(hash) == 0 || len(hash) > HashLen {
		return nil, errors.New("message hash is illegal")
	}
	s, err := sig.SerializeRSV()
	if err != nil {
		return nil, err
	}
	return ParsePublicKey(secp256k1.RecoverPubkey(hash[:], s))
}

// Verify verifies the signature of hash using the public key.
func (sig *Signature) Verify(msg []byte, pubKey *PublicKey) bool {
	if len(msg) == 0 || len(msg) > HashLen || pubKey == nil {
		return false
	}
	s, err := sig.SerializeRSV()
	if err != nil {
		return false
	}
	ret := secp256k1.VerifySignature(msg, s, pubKey.bytes)
	return ret != 0
}

// String returns the string representation.
func (sig *Signature) String() string {
	if sig == nil || len(sig.bytes) == 0 {
		return "[empty]"
	}
	if sig.hasV {
		return "0x" + hex.EncodeToString(sig.bytes)
	}
	return "0x" + hex.EncodeToString(sig.bytes[:SignatureLenRaw]) + "[no V]"
}
