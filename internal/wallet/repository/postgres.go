package postgres

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/pkg/database"

	"github.com/google/uuid"
)

type WalletRepository struct {
	db *database.Client
}

func NewWalletRepository(db *database.Client) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r WalletRepository) GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error) {
	sql := `SELECT id, balance FROM wallets WHERE id = $1`

	row := r.db.QueryRow(ctx, sql, walletId)

	wallet := &models.Wallet{}
	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r WalletRepository) UpdateWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	sql := `UPDATE wallets SET balance = $1 WHERE id = $2`

	err := r.db.Exec(ctx, sql, wallet.Balance, wallet.ID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
