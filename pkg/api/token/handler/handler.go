package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/token"
	"juno/pkg/api/token/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger       logrus.FieldLogger
	tokenService token.Service
}

func New(logger logrus.FieldLogger, tokenService token.Service) *Handler {
	return &Handler{
		logger:       logger,
		tokenService: tokenService,
	}
}

func (h *Handler) Balance(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	balance, err := h.tokenService.Balance(u.ID)

	if err != nil {
		c.JSON(500, dto.NewErrorBalanceResponse(err.Error()))
		return
	}

	c.JSON(200, dto.NewSuccessBalanceResponse(balance))
}

func (h *Handler) Deposit(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorBalanceResponse(err.Error()))
		return
	}

	if err := h.tokenService.Deposit(u.ID, req.Amount); err != nil {
		c.JSON(500, dto.NewErrorBalanceResponse(err.Error()))
		return
	}

	c.JSON(200, dto.NewSuccessDepositResponse())
}
