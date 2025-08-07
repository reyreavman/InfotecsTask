package wallet

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error)
}
