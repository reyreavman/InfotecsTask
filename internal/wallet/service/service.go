package service

import (
	"context"
	"errors"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Реализация сервиса
// Ответственна за обработку ошибок и маппинг моделей
type WalletService struct {
	walletRepository wallet.Repository
}

func NewWalletService(walletRepository wallet.Repository) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
	}
}

func (s WalletService) GetWallet(ctx context.Context, walletId uuid.UUID) (*models.Wallet, error) {
	walletToReturn, err := s.walletRepository.GetWallet(ctx, walletId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wallet.ErrWalletNotFound
		}
		return nil, err
	}
	return walletToReturn, nil
}