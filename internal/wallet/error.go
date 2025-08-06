package wallet

import "errors"

var ErrWalletNotFound = errors.New("Wallet not found")
var ErrWalletWithNegativeBalance = errors.New("Wallet has negative balance")
var ErrFailedToUpdateWallet = errors.New("Failed to update wallet")