package model

type Wallet struct {
	Id      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type WalletOperation struct {
	Id            string
	OperationType string
	Amount        float64
}
