package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/transaction"
	"juno/pkg/api/transaction/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	transactionService transaction.Service
}

func New(transactionService transaction.Service) *Handler {
	return &Handler{
		transactionService: transactionService,
	}
}

func (h *Handler) List(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	transactions, err := h.transactionService.GetTransactionsByUserID(u.ID)
	if err != nil {
		c.JSON(500, dto.NewErrorListResponse(err.Error()))
		return
	}
	c.JSON(200, dto.NewSuccessListResponse(transactions))
}
