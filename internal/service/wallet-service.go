package service

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"wallet-api/internal/config"
	"wallet-api/internal/dto"
	"wallet-api/internal/logger"
	"wallet-api/internal/mapper"
	"wallet-api/internal/repository"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type WalletService interface {
	GetWalletById(id uuid.UUID) (*dto.WalletDTO, error)
	OperationWithWalletByID(dto dto.WalletOperationRequestDTO) error
	EnqueueOperation(op dto.WalletOperationRequestDTO)
}

type walletServiceImpl struct {
	repos          repository.WalletRepository
	operationQueue chan dto.WalletOperationRequestDTO
	cache          sync.Map // Cache for wallet data
	cacheTTL       time.Duration
	log            *slog.Logger
}

type WalletOperationJob struct {
	Operation dto.WalletOperationRequestDTO
	Resp      chan error
}

func NewWalletServiceImpl(repos repository.WalletRepository, logger *slog.Logger) *walletServiceImpl {
	service := walletServiceImpl{
		repos:          repos,
		operationQueue: make(chan dto.WalletOperationRequestDTO, 1000),
		log:            logger,
	}

	return &service
}

func (s *walletServiceImpl) GetWalletById(id uuid.UUID) (*dto.WalletDTO, error) {
	op := "service.wallet-service.GetWalletById"

	s.log.Info("Start work in service", "op", op)

	if cachedWallet, found := s.cache.Load(id); found {
		s.log.Info("Cache hit for wallet", "op", op, "walletID", id)
		return cachedWallet.(*dto.WalletDTO), nil
	}

	model, err := s.repos.GetWalletById(id)
	if err != nil {
		s.log.Error("Error from repos", "op", op)
		return nil, err
	}

	dto := mapper.WalletModelToDto(*model)
	s.cache.Store(id, &dto)

	go s.expireCache(id, s.cacheTTL)

	return &dto, nil
}

func (s *walletServiceImpl) OperationWithWalletByID(dto dto.WalletOperationRequestDTO) error {
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

	err := s.repos.OperationWithWalletByID(model)

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

func (s *walletServiceImpl) EnqueueOperation(op dto.WalletOperationRequestDTO) {
	s.operationQueue <- op
}

func (s *walletServiceImpl) walletWorker(workerID int) {
	op := fmt.Sprintf("service.wallet-service.Worker-%d", workerID)

	for operation := range s.operationQueue {
		s.log.Info("Processing operation", "worker", workerID, "op", op)

		err := s.OperationWithWalletByID(operation)
		if err != nil {
			s.log.Error("Failed to process operation", "worker", workerID, "error", err, "op", op)
		} else {
			s.log.Info("Successfully processed operation", "worker", workerID, "op", op)
		}
	}
}

func (s *walletServiceImpl) StartWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go s.walletWorker(i)
	}
}

func (s *walletServiceImpl) expireCache(id uuid.UUID, ttl time.Duration) {
	time.Sleep(ttl)
	s.cache.Delete(id)
	s.log.Info("Cache expired for wallet", "walletID", id)
}
