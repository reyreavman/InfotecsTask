package facade

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
)

// Интерфейс, который должна реализовывать каждая структура, которая "хочет быть" фасадом =)
// Необходим для того, чтобы хендлер не думал к какому сервису бежать при получении того или иного запроса
// Предоставляет единый объект для работы со всей системой
type Facade interface {
	CreateTransaction(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error)
	GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error)
	GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error)
	GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error)
}
