package acala

import (
	"context"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/pack"
)

func (client *Client) Mint(ctx context.Context, minterKey signature.KeyringPair, phash, nhash pack.Bytes32, sig pack.Bytes65, amount uint64) (pack.Bytes, error) {
	opts := types.SerDeOptions{NoPalletIndices: true}
	types.SetSerDeOptions(opts)

	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get metadata: %v", err)
	}

	alice := types.NewAddressFromAccountID(minterKey.PublicKey)
	c, err := types.NewCall(meta, "RenVmBridge.mint", alice, phash, types.NewUCompactFromUInt(amount), nhash, sig)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("construct call: %v", err)
	}

	hash, err := client.api.RPC.Author.SubmitExtrinsic(types.NewExtrinsic(c))
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("submit extrinsic: %v", err)
	}

	return pack.NewBytes(hash[:]), nil
}

func (client *Client) Burn(ctx context.Context, burnerKey signature.KeyringPair, recipient [20]byte, amount uint64) (pack.Bytes, error) {
	opts := types.SerDeOptions{NoPalletIndices: false}
	types.SetSerDeOptions(opts)

	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get metadata: %v", err)
	}

	c, err := types.NewCall(meta, "RenVmBridge.burn", recipient, types.NewUCompactFromUInt(amount))
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("construct call: %v", err)
	}

	ext := types.NewExtrinsic(c)

	genesisHash, err := client.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get blockhash: %v", err)
	}

	rv, err := client.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get runtime version: %v", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", burnerKey.PublicKey, nil)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("create storage key: %v", err)
	}

	var accountInfo types.AccountInfo
	ok, err := client.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		return pack.Bytes{}, fmt.Errorf("get storage: %v", err)
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
		return pack.Bytes{}, fmt.Errorf("sign extrinsic: %v", err)
	}

	hash, err := client.api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("submit extrinsic: %v", err)
	}

	return pack.NewBytes(hash[:]), nil
}
