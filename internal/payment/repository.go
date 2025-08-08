package payment

import (
	"context"
	"infotecstechtask/internal/models"
)

// Интерфейс репозитория для создания транзакций
// Выделил операцию в отдельный интерфейс, чтобы все операции во время создания и выполнения транзакции выполнялись в одной БД транзакции
type Repository interface {
	CreatePayment(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error)
}
