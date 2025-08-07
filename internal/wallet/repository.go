package wallet

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
)

// Интерфейс репозитория для получения баланса кошелька
type Repository interface {
	GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error)
}
