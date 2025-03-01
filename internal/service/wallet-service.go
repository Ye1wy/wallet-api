package service

import (
	"context"
	"log/slog"
	"wallet-api/internal/dto"
	"wallet-api/internal/repository"
)

type WalletService interface {
	GetWalletById(ctx context.Context, id string) (dto.WalletDTO, error)
	DepositToWalletByID(ctx context.Context, dto dto.WalletOperationRequestDTO) error
	WithdrawFromWalletByID(ctx context.Context, dto dto.WalletOperationRequestDTO) error
}

type walletServiceImpl struct {
	repos repository.WalletRepository
	log   *slog.Logger
}

func NewWalletServiceImpl(repos repository.WalletRepository, logger *slog.Logger) *walletServiceImpl {
	return &walletServiceImpl{
		repos: repos,
		log:   logger,
	}
}

func (s *walletServiceImpl) GetWalletById(ctx context.Context, id string) (*dto.WalletDTO, error) {
	return nil, nil
}

func (s *walletServiceImpl) DepositToWalletByID(ctx context.Context, dto dto.WalletOperationRequestDTO) error {
	return nil
}

func (s *walletServiceImpl) WithdrawFromWalletByID(ctx context.Context, dto dto.WalletOperationRequestDTO) error {
	return nil
}
