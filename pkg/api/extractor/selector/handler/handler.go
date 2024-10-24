package handler

import (
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/api/extractor/selector/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service selector.Service
	policy  selector.Policy
}

func New(policy selector.Policy, service selector.Service) *Handler {
	return &Handler{
		policy:  policy,
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateSelectorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateSelectorResponse(err))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			sel, err := h.service.Create(u.ID, req.Name, req.Value, selector.Visibility(req.Visibility))

			if err != nil {
				c.JSON(500, dto.NewErrorCreateSelectorResponse(err))
				return
			}

			c.JSON(201, dto.NewSuccessCreateSelectorResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateSelectorResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorCreateSelectorResponse(err))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorGetSelectorResponse(err))
		return
	}

	sel, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorGetSelectorResponse(err))
		return
	}

	h.policy.CanRead(c.Request.Context(), sel).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetSelectorResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetSelectorResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorGetSelectorResponse(err))
		})
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	sels, err := h.service.ListByUserID(u.ID)

	if err != nil {
		c.JSON(500, dto.NewErrorListSelectorResponse(err))
		return
	}

	h.policy.CanList(c.Request.Context(), sels).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListSelectorResponse(sels))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListSelectorResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListSelectorResponse(err))
		})
}
