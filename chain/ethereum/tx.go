package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/evm"
	"github.com/renproject/pack"
)

// TxBuilder represents a transaction builder that builds transactions to be
// broadcasted to the ethereum network. The TxBuilder is configured using a
// chain id.
type TxBuilder struct {
	ChainID *big.Int
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(chainID *big.Int) TxBuilder {
	return TxBuilder{chainID}
}

// BuildTx receives transaction fields and constructs a new transaction.
func (txBuilder TxBuilder) BuildTx(ctx context.Context, from, to address.Address, value, nonce, gas, gasTipCap, gasFeeCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	toAddr, err := NewAddressFromHex(string(pack.String(to)))
	if err != nil {
		return nil, fmt.Errorf("bad to address '%v': %v", to, err)
	}
	addr := common.Address(toAddr)
	return &evm.Tx{
		EthTx: types.NewTx(&types.DynamicFeeTx{
			ChainID:   txBuilder.ChainID,
			Nonce:     nonce.Int().Uint64(),
			GasTipCap: gasTipCap.Int(),
			GasFeeCap: gasFeeCap.Int(),
			Gas:       gas.Int().Uint64(),
			To:        &addr,
			Value:     value.Int(),
			Data:      payload,
		}),
		Signer: types.LatestSignerForChainID(txBuilder.ChainID),
	}, nil
}
