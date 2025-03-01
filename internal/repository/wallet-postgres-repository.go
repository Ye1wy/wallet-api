package repository

import (
	"context"
	"log/slog"
	"wallet-api/internal/logger"
	"wallet-api/internal/model"

	"github.com/jackc/pgx/v5"
)

type WalletRepository interface {
	GetWalletById(ctx context.Context, id string) (*model.Wallet, error)
	DepositToWalletByID(ctx context.Context, model model.Wallet) error
	WithdrawFromWalletByID(ctx context.Context, model model.Wallet) error
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

func (r *postgresWalletRepository) GetWalletById(ctx context.Context, id string) (*model.Wallet, error) {
	op := "repository.wallet-postgres-repository.GetWalletById"

	query := "SELECT id, balance FROM wallet"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error("Cannot take data: ", logger.Err(err), "op", op)
		return nil, err
	}
	defer rows.Close()

	var wallet model.Wallet
	if err := rows.Scan(&wallet.Id, &wallet.Balance); err != nil {
		r.log.Error("Failed to convert data to model struct", "op", op)
		return nil, err
	}

	r.log.Info("Successfully geted wallet data by id from database", "op", op)
	return &wallet, nil
}

func (r *postgresWalletRepository) DepositToWalletByID(ctx context.Context, model model.Wallet) error {
	return nil
}

func (r *postgresWalletRepository) WithdrawFromWalletByID(ctx context.Context, model model.Wallet) error {
	return nil
}
