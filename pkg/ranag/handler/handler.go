package handler

import (
	"juno/pkg/ranag"
	"juno/pkg/ranag/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ranag.Service
}

func New(service ranag.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RangeAggregate(c *gin.Context) {

	var req dto.RangeAggregatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorRangeAggregatorResponse(err))
		return
	}

	res, err := h.service.RangeAggregate(req.Offset, req.Total, req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.NewSuccessRangeAggregatorResponse(res))
}
