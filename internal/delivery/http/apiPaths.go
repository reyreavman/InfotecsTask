package http

const (
	BASED_PATH = "/api"
	SWAGGER = "/swagger/*any"

	SEND = "/send"
	TRANSACTIONS = "/transactions"
	GET_WALLET_BALANCE = "/wallet/:walletId/balance"

	FULL_SEND = "/api/send"
	FULL_TRANSACTIONS = "/api/transactions"
	FULL_GET_WALLET_BALANCE = "/api/wallet/:walletId/balance"
)