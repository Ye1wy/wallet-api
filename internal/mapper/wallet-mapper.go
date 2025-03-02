package mapper

import (
	"wallet-api/internal/dto"
	"wallet-api/internal/model"
)

func WalletModelToDto(model model.Wallet) dto.WalletDTO {
	return dto.WalletDTO{
		Balance: model.Balance,
	}
}

func WalletDtoToModel(dto dto.WalletDTO) model.Wallet {
	return model.Wallet{
		Balance: dto.Balance,
	}
}

func WalletOperationModelToDto(model model.WalletOperation) dto.WalletOperationRequestDTO {
	return dto.WalletOperationRequestDTO{
		Id:            model.Id,
		OperationType: model.OperationType,
		Amount:        model.Amount,
	}
}

func WalletOperationDtoToModel(dto dto.WalletOperationRequestDTO) model.WalletOperation {
	return model.WalletOperation{
		Id:            dto.Id,
		OperationType: dto.OperationType,
		Amount:        dto.Amount,
	}
}
