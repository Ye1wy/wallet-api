package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"wallet-api/internal/config"
	"wallet-api/internal/logger"
	"wallet-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type WalletRepository interface {
	GetWalletById(ctx *gin.Context, id string) (*model.Wallet, error)
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

func (r *postgresWalletRepository) GetWalletById(ctx *gin.Context, id string) (*model.Wallet, error) {
	op := "repository.wallet-postgres-repository.GetWalletById"

	query := "SELECT id, balance FROM wallet"

	row, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error("Cannot take data: ", logger.Err(err), "op", op)
		return nil, err
	}
	defer row.Close()

	var wallet model.Wallet
	if err := row.Scan(&wallet.Id, &wallet.Balance); err != nil {
		r.log.Error("Failed to convert data to model struct", logger.Err(err), "op", op)
		return nil, err
	}

	return &wallet, nil
}

func (r *postgresWalletRepository) OperationWithWalletByID(ctx *gin.Context, model model.WalletOperation) error {
	op := "repository.wallet-postgres-repository.DepositToWalletByID"
	var request string

	switch model.OperationType {
	case config.OperationDeposit:
		request = "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	case config.OperationWithdraw:
		request = "UPDATE wallets SET balance = balance - $1 WHERE id = $2"
	default:
		r.log.Error("Somting is wrong", "op", op)
		return errors.New("operation type is invalid")
	}

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
