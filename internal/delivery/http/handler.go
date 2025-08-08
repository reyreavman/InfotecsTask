package http

import (
	"context"
	"errors"
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/internal/wallet"
	"net/http"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transaction, err := h.facade.CreateTransaction(ctx, createTransactionRequest)
	if err != nil {
		if errors.Is(err, payment.ErrSenderWalletNotFound) || errors.Is(err, payment.ErrRecipientWalletNotFound) || errors.Is(err, payment.ErrSenderAndRecipientSame) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				models.Error{
					Error: err.Error(),
				},
			)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	params := c.MustGet("validatedParams").(*models.GetTransactionWithCountRequest)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if params.Count != nil {
		transactions, err := h.facade.GetTransactions(ctx, *params.Count)
		if err != nil {
			if errors.Is(err, wallet.ErrWalletNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, transactions)
	} else {
		transactions, err := h.facade.GetAllTransactions(ctx)
		if err != nil {
			if errors.Is(err, wallet.ErrWalletNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				c.AbortWithStatus(http.StatusServiceUnavailable)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	walletToReturn, err := h.facade.GetWallet(ctx, walletId)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, walletToReturn)
}
