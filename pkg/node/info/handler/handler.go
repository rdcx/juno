package handler

import (
	"juno/pkg/node/info"
	"juno/pkg/node/info/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service info.Service
}

func New(service info.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Info(c *gin.Context) {
	info, err := h.service.GetInfo()

	if err != nil {
		c.JSON(500, dto.NewErrorInfoResponse(err.Error()))
		return
	}

	c.JSON(200, dto.NewSuccessInfoResponse(info))
}
