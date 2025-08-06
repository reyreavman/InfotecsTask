package facade

import (
	"context"
	"infotecstechtask/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFacade struct {
	mock.Mock
}

func (m *MockFacade) CreateTransaction(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error) {
	args := m.Called(ctx, createTransactionRequest)

	var resp *models.TransactionResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*models.TransactionResponse)
	}
	
	return resp, args.Error(1)
}

func (m *MockFacade) GetTransactions(ctx context.Context, count int) ([]*models.TransactionResponse, error) {
	args := m.Called(ctx, count)
	
	var resp []*models.TransactionResponse
	if args.Get(0) != nil {
		resp = args.Get(0).([]*models.TransactionResponse)
	}
	
	return resp, args.Error(1)
}

func (m *MockFacade) GetAllTransactions(ctx context.Context) ([]*models.TransactionResponse, error) {
	args := m.Called(ctx)
	
	var resp []*models.TransactionResponse
	if args.Get(0) != nil {
		resp = args.Get(0).([]*models.TransactionResponse)
	}
	
	return resp, args.Error(1)
}

func (m *MockFacade) GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error) {
	args := m.Called(ctx, walletId)
	
	var wallet *models.Wallet
	if args.Get(0) != nil {
		wallet = args.Get(0).(*models.Wallet)
	}
	
	return wallet, args.Error(1)
}