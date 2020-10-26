package acala

import "github.com/centrifuge/go-substrate-rpc-client/types"

type eventMinted struct {
	Phase  types.Phase
	Owner  types.AccountID
	Amount types.U128
	Topics []types.Hash
}

type eventBurnt struct {
	Phase  types.Phase
	Owner  types.AccountID
	Dest   types.Bytes
	Amount types.U128
	Topics []types.Hash
}

type eventDeposited struct {
	Phase      types.Phase
	CurrencyId [2]byte
	Who        types.AccountID
	Amount     types.U128
	Topics     []types.Hash
}

type eventWithdrawn struct {
	Phase      types.Phase
	CurrencyId [2]byte
	Who        types.AccountID
	Amount     types.U128
	Topics     []types.Hash
}

type eventTreasury struct {
	Phase   types.Phase
	Deposit types.U128
	Topics  []types.Hash
}

type RenVmBridgeEvents struct {
	types.EventRecords
	Currencies_Deposited  []eventDeposited
	RenVmBridge_Minted    []eventMinted
	Currencies_Withdrawn  []eventWithdrawn
	RenVmBridge_Burnt     []eventBurnt
	AcalaTreasury_Deposit []eventTreasury
}
