package solana

type AccountContext struct {
	Slot int `json:"slot"`
}

type AccountValue struct {
	Data       string `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int    `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int    `json:"rentEpoch"`
}

type ResponseGetAccountInfo struct {
	Context AccountContext `json:"context"`
	Value   AccountValue   `json:"value"`
}

type BurnLog struct {
	Amount    int      `json:"amount"`
	Recipient [25]byte `json:"recipient"`
}
