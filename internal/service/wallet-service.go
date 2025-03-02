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
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type WalletService interface {
	GetWalletById(ctx *gin.Context, id uuid.UUID) (*dto.WalletDTO, error)
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

func (s *walletServiceImpl) GetWalletById(ctx *gin.Context, id uuid.UUID) (*dto.WalletDTO, error) {
	op := "service.wallet-service.GetWalletById"

	s.log.Info("Start work in service", "op", op)

	model, err := s.repos.GetWalletById(ctx, id)
	if err != nil {
		s.log.Error("Error from repos", "op", op)
		return nil, err
	}

	s.log.Info("Start mapping from model to dto", "op", op)
	dto := mapper.WalletModelToDto(*model)
	s.log.Info("Mapping successfully done", "op", op)

	return &dto, nil
}

func (s *walletServiceImpl) OperationWithWalletByID(ctx *gin.Context, dto dto.WalletOperationRequestDTO) error {
	op := "service.wallet-service.OperationWithWalletByID"
	s.log.Info("Start validate amount", "op", op)

	if dto.Amount < 0 {
		s.log.Error("Amount is not valid", "op", op)
		return errors.New(config.AmountIsNotValid)
	}

	s.log.Info("Start validate operation", "op", op)

	if dto.OperationType != config.OperationDeposit && dto.OperationType != config.OperationWithdraw {
		s.log.Error("Operation is not valid", "op", op)
		return errors.New(config.InvalidOperation)
	}

	s.log.Info("Start convertation dto to model", "op", op)
	model := mapper.WalletOperationDtoToModel(dto)

	s.log.Info("Initialize operation", "op", op)

	err := s.repos.OperationWithWalletByID(ctx, model)

	if err != nil {
		if err == pgx.ErrNoRows {
			s.log.Warn("Wallet is not found", "op", op)
			return err
		}

		s.log.Error("Failed to run operation with wallet", logger.Err(err), "op", op)
		return err
	}

	return nil
}
