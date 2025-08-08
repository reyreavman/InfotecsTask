package wallet

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
)

// Интерфейс сервиса
// Содержит в себе метод для получения баланса кошелька
type Service interface {
	GetWallet(ctx context.Context, walletId uuid.UUID) (*models.WalletResponse, error)
}
