package transaction

import (
	"context"
	"infotecstechtask/internal/models"
)

// Интерфейс репозитория
// Содержит в себе методы для получения списка транзакций
type Repository interface {
	GetTransactions(ctx context.Context, count int) ([]*models.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]*models.Transaction, error)
}
