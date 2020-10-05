package icon

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/haltingstate/secp256k1-go"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/multichain/chain/icon/intconv"
	"github.com/renproject/multichain/chain/icon/transaction"
	"github.com/renproject/pack"
	"golang.org/x/crypto/sha3"
)

type Signature struct {
	bytes []byte // 65 bytes of [R|S|V]
	hasV  bool
}

type HexBytes pack.String

func (hs HexBytes) Bytes() []byte {
	bs, _ := hex.DecodeString(string(hs[2:]))
	return bs
}

type HexInt struct {
	big.Int
}

type HexUint16 struct {
	Value uint16
}

func (i HexUint16) String() string {
	return intconv.FormatInt(int64(i.Value))
}

type HexInt64 struct {
	Value int64
}

func (i HexInt64) String() string {
	return intconv.FormatInt(i.Value)
}

type RawMessage []byte

// PrivateKey is a type representing a private key.
// for both private key and public key
type PrivateKey pack.Bytes

type Tx struct {
	txHash    pack.Bytes
	Version   HexUint16
	from      address.Address
	to        address.Address
	value     *HexInt
	StepLimit HexInt
	TimeStamp HexInt64
	NID       *HexInt
	nonce     *HexInt
	Signature Signature
	Data      RawMessage
	DataType  *string
}

func (tx *Tx) Hash() pack.Bytes {
	h, err := tx.Sighashes()
	if err != nil {
		tx.txHash = pack.Bytes{}
	} else {
		b := h[0][:]
		tx.txHash = pack.NewBytes(b)
	}
	return tx.txHash
}

// From returns the sender of the transaction
func (tx Tx) From() address.Address {
	return tx.from
}

// To returns the recipients of the transaction.
func (tx Tx) To() address.Address {
	return tx.to
}

// Value returns the values being transferred in a transaction. For icon
// chain, there can be multiple messages (each with a different value being
// transferred) in a single transaction.
func (tx *Tx) Value() pack.U256 {
	u, _ := strconv.ParseUint(tx.value.String()[2:], 16, 64)
	return pack.NewU256FromU64(pack.U64(u))
}

// Nonce returns the transaction count of the transaction sender.
func (tx Tx) Nonce() pack.U256 {
	u, _ := strconv.ParseUint(tx.nonce.String()[2:], 16, 64)
	return pack.NewU256FromU64(pack.U64(u))
}

// Payload returns the memo attached to the transaction.
func (tx *Tx) Payload() contract.CallData {
	return contract.CallData(pack.Bytes(make([]byte, 0)))
}

// Sighashes that need to be signed before this transaction can be submitted.
func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	sighashes := make([]pack.Bytes32, 0)
	sha := bytes.NewBuffer(nil)
	sha.Write([]byte("icx_sendTransaction"))

	// data
	if tx.Data != nil {
		sha.Write([]byte(".data."))
		if len(tx.Data) > 0 {
			var obj interface{}
			if err := json.Unmarshal(tx.Data, &obj); err != nil {
				return nil, err
			}
			bs, err := transaction.SerializeValue(obj)
			if err != nil {
				return nil, err
			}
			sha.Write(bs)
		}
	}

	// dataType
	if tx.DataType != nil {
		sha.Write([]byte(".dataType."))
		sha.Write([]byte(*tx.DataType))
	}

	// from
	sha.Write([]byte(".from."))
	sha.Write([]byte(tx.from))

	// nid
	sha.Write([]byte(".nid."))
	sha.Write([]byte(tx.NID.String()))

	// nonce
	sha.Write([]byte(".nonce."))
	sha.Write([]byte(tx.nonce.String()))

	// stepLimit
	sha.Write([]byte(".stepLimit."))
	sha.Write([]byte(tx.StepLimit.String()))

	// timestamp
	sha.Write([]byte(".timestamp."))
	sha.Write([]byte(tx.TimeStamp.String()))

	// to
	sha.Write([]byte(".to."))
	sha.Write([]byte(tx.To()))

	// value
	sha.Write([]byte(".value."))
	sha.Write([]byte(tx.value.String()))

	// version
	sha.Write([]byte(".version."))
	sha.Write([]byte(tx.Version.String()))
	d := sha3.Sum256(sha.Bytes())
	hash := d[:]
	sighash := [32]byte{}
	copy(sighash[:], hash)
	sighashes[0] = pack.NewBytes32(sighash)
	return sighashes, nil
}

var globalLock sync.Mutex

const (
	// SignatureLenRawWithV is the bytes length of signature including V value
	SignatureLenRawWithV = 65
	// SignatureLenRaw is the bytes length of signature not including V value
	SignatureLenRaw = 64
	// HashLen is the bytes length of hash for signature
	HashLen = 32
)

// NewSignature calculates an ECDSA signature including V, which is 0 or 1.
func NewSignature(hash []byte, privKey PrivateKey) (*Signature, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if len(hash) == 0 || len(hash) > HashLen || privKey == nil {
		return nil, errors.New("Invalid arguments")
	}
	return ParseSignature(secp256k1.Sign(hash, privKey))
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

// Sign ...
func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	for _, rsv := range signatures {
		data := rsv[:]
		pk := PrivateKey(pubKey)
		sig, err := NewSignature(data, pk)
		if err != nil {
			return err
		}
		tx.Signature = *sig
	}
	return nil
}

var txSerializeExcludes = map[string]bool{"signature": true}

// Serialize the transaction.
func (tx *Tx) Serialize() (pack.Bytes, error) {
	tx.TimeStamp = HexInt64{time.Now().UnixNano() / int64(time.Microsecond)}
	js, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return nil, err
	}
	return pack.Bytes(bs), nil
}

type TxBuilder struct{}

// BuildTx consumes a list of MsgSend to build and return a transaction.
// This transaction is unsigned, and must be signed before submitting to the chain.
func (txBuilder TxBuilder) BuildTx(from address.Address, to address.Address, value *HexInt, nonce *HexInt, Version HexUint16, txHash pack.Bytes, StepLimit HexInt, TimeStamp HexInt64, NID *HexInt) (*Tx, error) {
	if len(from) != 42 {
		//return nil, fmt.Errorf("Address sending is invalid", from)
	}

	if len(to) != 42 {
		//return nil, fmt.Errorf("Address receiving is invalid", to)
	}

	return &Tx{
		from:      from,
		to:        to,
		value:     value,
		nonce:     nonce,
		Version:   Version,
		StepLimit: StepLimit,
		TimeStamp: TimeStamp,
		NID:       NID,
	}, nil
}
