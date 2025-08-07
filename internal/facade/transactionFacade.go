package facade

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/internal/transaction"
	"infotecstechtask/internal/wallet"

	"github.com/google/uuid"
)

// Реализация интерфейса Facade
// Содержит в себе WalletService, TransactionService и PaymentRepository
type TransactionFacade struct {
	walletService      wallet.Service
	transactionService transaction.Service
	paymentRepository  payment.Repository
}

func NewFacade(walletService wallet.Service, transactionService transaction.Service, paymentRepository payment.Repository) *TransactionFacade {
	return &TransactionFacade{
		walletService:      walletService,
		transactionService: transactionService,
		paymentRepository:  paymentRepository,
	}
}

func (f TransactionFacade) CreateTransaction(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error) {
	transaction, err := f.paymentRepository.CreatePayment(ctx, createTransactionRequest)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (f TransactionFacade) GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error) {
	return f.transactionService.GetTransactions(ctx, count)
}

func (f TransactionFacade) GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error) {
	return f.transactionService.GetAllTransactions(ctx)
}

func (f TransactionFacade) GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error) {
	return f.walletService.GetWallet(ctx, walletId)
}
