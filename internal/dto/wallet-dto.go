package dto

import "github.com/gofrs/uuid/v5"

type WalletDTO struct {
	Balance float64 `json:"balance"`
}

type WalletOperationRequestDTO struct {
	Id            uuid.UUID `json:"id" binding:"required"`
	OperationType string    `json:"operation_type" binding:"required"`
	Amount        float64   `json:"amount" binding:"required"`
}

type ErrorDTO struct {
	Error      string
	StatusCode int
}
