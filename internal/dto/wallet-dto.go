package dto

type WalletDTO struct {
	Balance float64 `json:"balance"`
}

type WalletOperationRequestDTO struct {
	Id            string  `json:"id" binding:"required"`
	OperationType string  `json:"operation_type" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}
