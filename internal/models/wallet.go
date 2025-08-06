package models

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID
	Balance float32
}

type GetWalletBalanceRequest struct {
	ID string `uri:"walletId" validate:"required,uuid"`
}