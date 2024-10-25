package handler

import (
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/field/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service field.Service
	policy  field.Policy
}

func New(policy field.Policy, service field.Service) *Handler {
	return &Handler{
		policy:  policy,
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateFieldRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateFieldResponse(err))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			sel, err := h.service.Create(u.ID, uuid.MustParse(req.SelectorID), req.Name, field.FieldType(req.Type))

			if err != nil {
				c.JSON(500, dto.NewErrorCreateFieldResponse(err))
				return
			}

			c.JSON(201, dto.NewSuccessCreateFieldResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorCreateFieldResponse(err))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorGetFieldResponse(err))
		return
	}

	sel, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorGetFieldResponse(err))
		return
	}

	h.policy.CanRead(c.Request.Context(), sel).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetFieldResponse(sel))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorGetFieldResponse(err))
		})
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	sels, err := h.service.ListByUserID(u.ID)

	if err != nil {
		c.JSON(500, dto.NewErrorListFieldResponse(err))
		return
	}

	h.policy.CanList(c.Request.Context(), sels).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListFieldResponse(sels))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListFieldResponse(err))
		})
}

func (h *Handler) ListBySelectorID(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	sels, err := h.service.ListBySelectorID(u.ID)

	if err != nil {
		c.JSON(400, dto.NewErrorListFieldResponse(err))
		return
	}

	h.policy.CanList(c.Request.Context(), sels).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListFieldResponse(sels))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListFieldResponse(err))
		})
}
