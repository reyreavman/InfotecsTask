package models

import (
	"github.com/google/uuid"
)

// Модель кошелька для клиента
type WalletResponse struct {
	ID      uuid.UUID
	Balance float64
}

// Модель кошелька, хранящаяся в БД
type Wallet struct {
	ID uuid.UUID
	Balance int
}

// Модель аккумулирующая в себе параметры запроса для получения баланса кошелька
type GetWalletBalanceRequest struct {
	ID string `uri:"walletId" validate:"required,uuid"`
}

func ToWalletResponse(wallet *Wallet) *WalletResponse {
	return &WalletResponse{
		ID:      wallet.ID,
		Balance: float64(wallet.Balance) / 100.0,
	}
}
