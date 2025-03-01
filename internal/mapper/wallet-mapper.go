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
