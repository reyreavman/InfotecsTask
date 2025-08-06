package service

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/wallet"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepository wallet.Repository
}

func NewWalletService(walletRepository wallet.Repository) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
	}
}

func (s WalletService) GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error) {
	return s.walletRepository.GetWallet(ctx, walletId)
}

func (s WalletService) UpdateWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	if err := checkNegativeBalance(wallet); err != nil {
		return nil, err
	}
	
	return s.walletRepository.UpdateWallet(ctx, wallet)
}

func checkNegativeBalance(w *models.Wallet) error {
	if w.Balance < 0 {
		return wallet.ErrWalletWithNegativeBalance
	}
	return nil
}
