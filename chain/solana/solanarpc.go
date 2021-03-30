package solana

// AccountContext is the JSON-interface of the account's context representing
// what slot the account's value has been returned for.
type AccountContext struct {
	Slot int `json:"slot"`
}

// AccountValue is the JSON-interface of the account's information.
type AccountValue struct {
	Data       string `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int    `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int    `json:"rentEpoch"`
}

// ResponseGetAccountInfo is the JSON-interface of the response for the
// getAccountInfo query.
type ResponseGetAccountInfo struct {
	Context AccountContext `json:"context"`
	Value   AccountValue   `json:"value"`
}

// BurnLog is the data stored in a burn log account, that is received in its
// Base58 encoded format as a part of the getAccountInfo response.
type BurnLog struct {
	Amount    int      `json:"amount"`
	Recipient [25]byte `json:"recipient"`
}

type Bytes32 = [32]byte

type gatewayRegistry struct {
	IsInitialised bool
	Owner         Bytes32
	Count         uint64
	Selectors     []Bytes32
	Gateways      []Bytes32
}
