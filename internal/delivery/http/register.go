package http

import (
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/middleware"
	"infotecstechtask/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Функция для регистрации эндпоинтов и соответствующих функций хендлера
func RegisterHTTPEndpoints(router *gin.RouterGroup, facade facade.Facade, validate *validator.Validate) {
	h := NewHandler(facade)

	api := router.Group(BASED_PATH)
	{
		api.POST(SEND, middleware.JSONValidation(models.CreateTransactionRequest{}, validate), h.CreateTransaction)
		api.GET(TRANSACTIONS, middleware.ParamsValidation(models.GetTransactionWithCountRequest{}, validate), h.GetTransactions)
		api.GET(GET_WALLET_BALANCE, middleware.ParamsValidation(models.GetWalletBalanceRequest{}, validate), h.GetWallet)
	}
}
