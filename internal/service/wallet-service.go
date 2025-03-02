package service

import (
	"errors"
	"log/slog"
	"wallet-api/internal/config"
	"wallet-api/internal/dto"
	"wallet-api/internal/logger"
	"wallet-api/internal/mapper"
	"wallet-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletService interface {
	GetWalletById(ctx *gin.Context, id string) (*dto.WalletDTO, error)
	OperationWithWalletByID(ctx *gin.Context, dto dto.WalletOperationRequestDTO) error
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

func (s *walletServiceImpl) GetWalletById(ctx *gin.Context, id string) (*dto.WalletDTO, error) {
	op := "service.wallet-service.GetWalletById"

	if err := uuid.Validate(id); err != nil {
		s.log.Error("id is not valid", "op", op)
		return nil, errors.New(config.UUIDIsNotValid)
	}

	model, err := s.repos.GetWalletById(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := mapper.WalletModelToDto(*model)

	return &dto, nil
}

func (s *walletServiceImpl) OperationWithWalletByID(ctx *gin.Context, dto dto.WalletOperationRequestDTO) error {
	op := "service.wallet-service.OperationWithWalletByID"

	if err := uuid.Validate(dto.Id); err != nil {
		s.log.Error("Id is not valid", "op", op)
		return errors.New(config.UUIDIsNotValid)
	}

	if dto.Amount < 0 {
		s.log.Error("Amount is not valid", "op", op)
		return errors.New(config.AmountIsNotValid)
	}

	if dto.OperationType != config.OperationDeposit && dto.OperationType != config.OperationWithdraw {
		s.log.Error("Operation id is not valid", "op", op)
		return errors.New(config.InvalidOperation)
	}

	model := mapper.WalletOperationDtoToModel(dto)

	err := s.repos.OperationWithWalletByID(ctx, model)
	if err == pgx.ErrNoRows {
		s.log.Warn("Wallet is not found", "op", op)
		return errors.New(config.UUIDIsNotValid)
	}

	if err != nil {
		s.log.Error("Failed to run operation with wallet", logger.Err(err), "op", op)
		return err
	}

	return nil
}
