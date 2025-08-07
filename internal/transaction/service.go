package transaction

import (
	"context"
	"infotecstechtask/internal/models"
)

// Интерфейс сервиса
// Содержит в себе методы для получения списка транзакций
type Service interface {
	GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error)
	GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error)
}
