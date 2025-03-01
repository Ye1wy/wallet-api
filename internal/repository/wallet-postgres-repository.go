package repository

import (
	"context"
	"log/slog"
	"wallet-api/internal/model"

	"github.com/jackc/pgx"
)

type WalletRepository interface {
	GetWalletById(ctx context.Context, id string) (*model.Wallet, error)
	DepositToWalletByID(ctx context.Context, model model.Wallet) error
	WithdrawFromWalletByID(ctx context.Context, model model.Wallet) error
}

type postgresWalletRepository struct {
	db     *pgx.Conn
	logger *slog.Logger
}

func NewPostgresWalletRepository(conn *pgx.Conn, log *slog.Logger) *postgresWalletRepository {
	return &postgresWalletRepository{
		db:     conn,
		logger: log,
	}
}

func (r *postgresWalletRepository) GetWalletById(ctx context.Context, id string) (*model.Wallet, error) {
	return nil, nil
}

func (r *postgresWalletRepository) DepositToWalletByID(ctx context.Context, model model.Wallet) error {
	return nil
}

func (r *postgresWalletRepository) WithdrawFromWalletByID(ctx context.Context, model model.Wallet) error {
	return nil
}
