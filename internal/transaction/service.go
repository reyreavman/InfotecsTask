package transaction

import (
	"context"
	"infotecstechtask/internal/models"
)

type Service interface {
	GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error)
	GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error)
}
