package handler

import (
	"juno/pkg/node/extraction"
	"juno/pkg/node/extraction/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger            *logrus.Logger
	extractionService extraction.Service
}

func New(
	logger *logrus.Logger,
	extracionService extraction.Service,
) *Handler {
	return &Handler{
		logger:            logger,
		extractionService: extracionService,
	}
}

func (h *Handler) Extract(c *gin.Context) {

	var req dto.ExtractionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorExtractionResponse(err))
		return
	}

	data, err := h.extractionService.Extract(req)

	if err != nil {
		h.logger.WithError(err).Error("failed to get titles")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessExtractionResponse(data))
}
