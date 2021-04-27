package solana

// AccountContext is the JSON-interface of the account's context representing
// what slot the account's value has been returned for.
type AccountContext struct {
	Slot int `json:"slot"`
}

// AccountValue is the JSON-interface of the account's information.
type AccountValue struct {
	Data       [2]string `json:"data"`
	Executable bool      `json:"executable"`
	Lamports   int       `json:"lamports"`
	Owner      string    `json:"owner"`
	RentEpoch  int       `json:"rentEpoch"`
}

// ResponseGetAccountInfo is the JSON-interface of the response for the
// getAccountInfo query.
type ResponseGetAccountInfo struct {
	Context AccountContext `json:"context"`
	Value   AccountValue   `json:"value"`
}
