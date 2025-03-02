package model

import "github.com/gofrs/uuid/v5"

type Wallet struct {
	Id      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}

type WalletOperation struct {
	Id            uuid.UUID
	OperationType string
	Amount        float64
}
