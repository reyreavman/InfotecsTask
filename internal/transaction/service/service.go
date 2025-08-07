package service

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/transaction"
)

// Реализация сервиса
// Ответственна за обработку ошибок и маппинг моделей
type TransactionService struct {
	transactionRepository transaction.Repository
}

func NewTransactionService(transactionRepository transaction.Repository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (s TransactionService) GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error) {
	transactions, err := s.transactionRepository.GetTransactions(ctx, count)

	return models.ToTransactionResponses(transactions), err
}

func (s TransactionService) GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error) {
	transactions, err := s.transactionRepository.GetAllTransactions(ctx)

	return models.ToTransactionResponses(transactions), err
}