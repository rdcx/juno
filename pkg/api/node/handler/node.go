package handler

import (
	"juno/pkg/api/node"
	"juno/pkg/api/node/dto"
	"juno/pkg/api/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger      *logrus.Logger
	nodeService node.Service
}

func New(l *logrus.Logger, ns node.Service) *Handler {
	return &Handler{
		logger:      l,
		nodeService: ns,
	}
}

func (h *Handler) Get(c *gin.Context) {
	u := c.MustGet("user").(*user.User)

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.nodeService.Get(u, id)

	if err != nil {
		h.logger.Error(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
				"node":  id,
			})
		c.JSON(404, dto.NewErrorGetNodeResponse(
			node.ErrNotFound.Error(),
		))
		return
	}

	c.JSON(200, dto.NewSuccessGetNodeResponse(n))
}

func (h *Handler) Create(c *gin.Context) {
	u := c.MustGet("user").(*user.User)

	var req dto.CreateNodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.nodeService.Create(u, req.Address, req.Shards)

	if err != nil {
		h.logger.Error(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
				"req":   req,
			})
		c.JSON(400, dto.NewErrorGetNodeResponse(
			err.Error(),
		))
		return
	}

	c.JSON(201, dto.NewSuccessGetNodeResponse(n))
}
