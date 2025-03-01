package dto

type WalletDTO struct {
	Balance float64 `json:"balance"`
}

type WalletOperationRequestDTO struct {
	WalletId      string  `json:"id" validate:"required,uuid"`
	OperationType string  `json:"operation_type" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64 `json:"amount" validate:"required,gte=0"`
}
