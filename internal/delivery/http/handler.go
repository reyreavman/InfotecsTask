package http

import (
	"errors"
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/internal/wallet"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Хендлер вызывает функции фасада и в зависимости от возвращаемых значений собирает ответ для клиента
// Запрос сюда попадает после прохожождения всех миддлваров
type Handler struct {
	facade facade.Facade
}

func NewHandler(facade facade.Facade) *Handler {
	return &Handler{
		facade: facade,
	}
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	createTransactionRequest := c.MustGet("validatedBody").(*models.CreateTransactionRequest)

	transaction, err := h.facade.CreateTransaction(c.Request.Context(), createTransactionRequest)
	if err != nil {
		if errors.Is(err, payment.ErrSenderWalletNotFound) || errors.Is(err, payment.ErrRecipientWalletNotFound) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest, 
				models.Error{
					Error: err.Error(),
				},
			)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	params := c.MustGet("validatedParams").(*models.GetTransactionWithCountRequest)

	if params.Count != nil {
		transactions, err := h.facade.GetTransactions(c.Request.Context(), *params.Count)
		if err != nil {
			if errors.Is(err, wallet.ErrWalletNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, transactions)
	} else {
		transactions, err := h.facade.GetAllTransactions(c.Request.Context())
		if err != nil {
			if errors.Is(err, wallet.ErrWalletNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, transactions)
	}
}

func (h *Handler) GetWallet(c *gin.Context) {
	walletId := uuid.MustParse(c.MustGet("validatedParams").(*models.GetWalletBalanceRequest).ID)

	walletToReturn, err := h.facade.GetWallet(c.Request.Context(), walletId)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, walletToReturn)
}
