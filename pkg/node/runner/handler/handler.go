package handler

import (
	"juno/pkg/node/runner"
	"juno/pkg/node/runner/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger        *logrus.Logger
	runnerService runner.Service
}

func New(
	logger *logrus.Logger,
	runnerService runner.Service,
) *Handler {
	return &Handler{
		logger:        logger,
		runnerService: runnerService,
	}
}

func (h *Handler) Titles(c *gin.Context) {
	titles, err := h.runnerService.Titles()

	if err != nil {
		h.logger.WithError(err).Error("failed to get titles")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, titles)
}

func (h *Handler) Execute(c *gin.Context) {
	var req dto.ExecuteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.runnerService.Execute(req.Src)

	if err != nil {
		h.logger.WithError(err).Error("failed to execute")
		c.JSON(http.StatusInternalServerError, dto.NewErrorExecuteResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessExecuteResponse(data))
}
