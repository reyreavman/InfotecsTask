package payment

import (
	"context"
	"infotecstechtask/internal/models"
)

type Repository interface {
	CreatePayment(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error)
}
