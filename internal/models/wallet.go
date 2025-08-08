package models

import "github.com/google/uuid"

// Модель кошелька
// Является полным отражением модели, которая хранится в БД
type Wallet struct {
	ID      uuid.UUID
	Balance float32
}

// Модель аккумулирующая в себе параметры запроса для получения баланса кошелька
type GetWalletBalanceRequest struct {
	ID string `uri:"walletId" validate:"required,uuid"`
}
