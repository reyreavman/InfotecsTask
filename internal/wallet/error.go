package wallet

import "errors"

// Список возможных ошибок бизнес-логики кошельков 
var ErrWalletNotFound = errors.New("Wallet not found")