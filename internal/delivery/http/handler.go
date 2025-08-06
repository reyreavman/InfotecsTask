package http

import (
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	params := c.MustGet("validatedParams").(*models.GetTransactionWithCountRequest)

	if params.Count != nil {
		transactions, err := h.facade.GetTransactions(c.Request.Context(), *params.Count)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, transactions)
	} else {
		transactions, err := h.facade.GetAllTransactions(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, transactions)
	}
}

func (h *Handler) GetWallet(c *gin.Context) {
	walletId := uuid.MustParse(c.MustGet("validatedParams").(*models.GetWalletBalanceRequest).ID)

	wallet, err := h.facade.GetWallet(c.Request.Context(), walletId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, wallet)
}

/*
 a35463b4-b74b-4467-8bc0-de84f2637fb0 |  100.00
 aecead58-d370-4d61-8e5a-2a086ff7d82e |  100.00
 5a97846b-f0b3-49cd-9dc9-4d9109c28e3c |  100.00
 7657ef49-dad0-4a12-8463-69a80e3cb8b6 |  100.00
 afed6c0c-09a4-42ce-a52e-cec33515d097 |  100.00
 427dc42f-5cb9-4e1f-8fda-e92420f86ac5 |  100.00
 bfbd92d5-e057-4df9-8da0-56c8fbd3c2e5 |  100.00
 3932945e-8682-4e12-9c77-c64b9a4edced |  100.00
 f15960ac-4f62-436b-b757-eccf8c24e4e5 |  100.00
 16800534-4a16-49cb-a85f-4c007e3a4b7a |  100.00
*/
