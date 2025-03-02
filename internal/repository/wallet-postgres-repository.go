package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"wallet-api/internal/config"
	"wallet-api/internal/logger"
	"wallet-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type WalletRepository interface {
	GetWalletById(ctx *gin.Context, id uuid.UUID) (*model.Wallet, error)
	OperationWithWalletByID(ctx *gin.Context, model model.WalletOperation) error
}

type postgresWalletRepository struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewPostgresWalletRepository(conn *pgx.Conn, logger *slog.Logger) *postgresWalletRepository {
	return &postgresWalletRepository{
		db:  conn,
		log: logger,
	}
}

func (r *postgresWalletRepository) GetWalletById(ctx *gin.Context, id uuid.UUID) (*model.Wallet, error) {
	op := "repository.wallet-postgres-repository.GetWalletById"

	query := "SELECT id, balance FROM wallets WHERE id = $1"

	r.log.Info("Start extraction data from data base", "op", op)

	var wallet model.Wallet
	err := r.db.QueryRow(ctx, query, id).Scan(&wallet.Id, &wallet.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.log.Error("Wallet not found", "op", op)
			return nil, pgx.ErrNoRows
		}
		r.log.Error("Failed to scan data", logger.Err(err), "op", op)
		return nil, err
	}

	r.log.Info("Extraction data is successfully done", "op", op)
	return &wallet, nil
}

func (r *postgresWalletRepository) OperationWithWalletByID(ctx *gin.Context, model model.WalletOperation) error {
	op := "repository.wallet-postgres-repository.DepositToWalletByID"
	var request string

	r.log.Info("Determinating operation type", "op", op)

	switch model.OperationType {
	case config.OperationDeposit:
		request = "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	case config.OperationWithdraw:
		request = "UPDATE wallets SET balance = balance - $1 WHERE id = $2"
	default:
		r.log.Error("Somting is wrong", "op", op)
		return errors.New("operation type is invalid")
	}

	r.log.Info("Updating wallet data", "op", op)

	res, err := r.db.Exec(ctx, request, model.Amount, model.Id)
	if err != nil {
		r.log.Error("Error in exec update request to data base", "op", op)
		return err
	}

	if res.RowsAffected() == 0 {
		r.log.Warn("Item is not found", "op", op)
		return fmt.Errorf("wallet not found")
	}

	return nil
}
