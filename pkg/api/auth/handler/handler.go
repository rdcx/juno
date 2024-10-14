package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger      *logrus.Logger
	authService auth.Service
}

func New(logger *logrus.Logger, authService auth.Service) *Handler {
	return &Handler{
		logger:      logger,
		authService: authService,
	}
}

func (h *Handler) Token(c *gin.Context) {
	var req dto.TokenRequest

	if err := c.BindJSON(&req); err != nil {
		h.logger.Error(err)
		c.JSON(400, dto.NewErrorTokenResponse(err.Error()))
		return
	}

	token, err := h.authService.Authenticate(req.Email, req.Password)

	if err != nil {
		h.logger.Error(err)
		c.JSON(400, dto.NewErrorTokenResponse(err.Error()))
		return
	}

	c.JSON(200, dto.NewSuccessTokenResponse(token))
}
