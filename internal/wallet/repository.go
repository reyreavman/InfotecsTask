package wallet

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error)
}
