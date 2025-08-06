package transaction

import (
	"context"
	"infotecstechtask/internal/models"
)

type Repository interface {
	GetTransactions(ctx context.Context, count int) ([]*models.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]*models.Transaction, error)
}
