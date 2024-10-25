package handler

import (
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/filter/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service filter.Service
	policy  filter.Policy
}

func New(policy filter.Policy, service filter.Service) *Handler {
	return &Handler{
		policy:  policy,
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateFilterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateFilterResponse(err))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			sel, err := h.service.Create(u.ID, req.Name, filter.FilterType(req.Type), req.Value)

			if err != nil {
				c.JSON(500, dto.NewErrorCreateFilterResponse(err))
				return
			}

			c.JSON(201, dto.NewSuccessCreateFilterResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateFilterResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorCreateFilterResponse(err))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorGetFilterResponse(err))
		return
	}

	sel, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorGetFilterResponse(err))
		return
	}

	h.policy.CanRead(c.Request.Context(), sel).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetFilterResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetFilterResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorGetFilterResponse(err))
		})
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	sels, err := h.service.ListByUserID(u.ID)

	if err != nil {
		c.JSON(500, dto.NewErrorListFilterResponse(err))
		return
	}

	h.policy.CanList(c.Request.Context(), sels).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListFilterResponse(sels))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListFilterResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListFilterResponse(err))
		})
}
