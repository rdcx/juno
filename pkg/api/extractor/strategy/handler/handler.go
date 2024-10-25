package handler

import (
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/extractor/strategy/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service strategy.Service
	policy  strategy.Policy
}

func New(policy strategy.Policy, service strategy.Service) *Handler {
	return &Handler{
		policy:  policy,
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateStrategyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateStrategyResponse(err))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			sel, err := h.service.Create(u.ID, req.Name)

			if err != nil {
				c.JSON(500, dto.NewErrorCreateStrategyResponse(err))
				return
			}

			c.JSON(201, dto.NewSuccessCreateStrategyResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateStrategyResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorCreateStrategyResponse(err))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorGetStrategyResponse(err))
		return
	}

	sel, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorGetStrategyResponse(err))
		return
	}

	h.policy.CanRead(c.Request.Context(), sel).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetStrategyResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetStrategyResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorGetStrategyResponse(err))
		})
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	sels, err := h.service.ListByUserID(u.ID)

	if err != nil {
		c.JSON(500, dto.NewErrorListStrategyResponse(err))
		return
	}

	h.policy.CanList(c.Request.Context(), sels).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListStrategyResponse(sels))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListStrategyResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListStrategyResponse(err))
		})
}
