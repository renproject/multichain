package acala

import (
	"fmt"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/pack"
)

func (client *Client) Mint(minterKey signature.KeyringPair, phash, nhash pack.Bytes32, sig pack.Bytes65, amount uint64) (pack.Bytes32, error) {
	opts := types.SerDeOptions{NoPalletIndices: true}
	types.SetSerDeOptions(opts)

	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("get metadata: %v", err)
	}

	alice := types.NewAddressFromAccountID(minterKey.PublicKey)
	c, err := types.NewCall(meta, "RenVmBridge.mint", alice, phash, types.NewUCompactFromUInt(amount), nhash, sig)
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("construct call: %v", err)
	}

	hash, err := client.api.RPC.Author.SubmitExtrinsic(types.NewExtrinsic(c))
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("submit extrinsic: %v", err)
	}

	return pack.NewBytes32(hash), nil
}

func (client *Client) Burn(burnerKey signature.KeyringPair, recipient pack.Bytes, amount uint64) (pack.Bytes32, error) {
	opts := types.SerDeOptions{NoPalletIndices: false}
	types.SetSerDeOptions(opts)

	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("get metadata: %v", err)
	}

	c, err := types.NewCall(meta, "RenVmBridge.burn", types.Bytes(recipient), types.NewUCompactFromUInt(amount))
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("construct call: %v", err)
	}

	ext := types.NewExtrinsic(c)

	genesisHash, err := client.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("get blockhash: %v", err)
	}

	rv, err := client.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("get runtime version: %v", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", burnerKey.PublicKey, nil)
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("create storage key: %v", err)
	}

	var accountInfo types.AccountInfo
	ok, err := client.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		return pack.Bytes32{}, fmt.Errorf("get storage: %v", err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(burnerKey, o)
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("sign extrinsic: %v", err)
	}

	sub, err := client.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return pack.Bytes32{}, fmt.Errorf("submit extrinsic: %v", err)
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				return pack.NewBytes32(status.AsInBlock), nil
			}
		case <-timeout:
			return pack.Bytes32{}, fmt.Errorf("timeout on tx confirmation")
		}
	}
}

type TokenAccount struct {
	Free     types.U128
	Reserved types.U128
	Frozen   types.U128
}

func (client *Client) Balance(user signature.KeyringPair) (pack.U256, error) {
	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.U256{}, fmt.Errorf("get metadata: %v", err)
	}

	key, err := types.CreateStorageKey(meta, "Tokens", "Accounts", user.PublicKey, []byte{0, 5})
	if err != nil {
		return pack.U256{}, fmt.Errorf("create storage key: %v", err)
	}

	var data TokenAccount
	ok, err := client.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil || !ok {
		return pack.U256{}, fmt.Errorf("get storage: %v", err)
	}

	balance := big.NewInt(0).SetBytes(data.Free.Bytes())
	return pack.NewU256FromInt(balance), nil
}
