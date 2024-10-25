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

func (h *Handler) AddSelector(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorAddSelectorResponse(err))
		return
	}

	var req dto.AddSelectorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorAddSelectorResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorAddSelectorResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.AddSelector(id, uuid.MustParse(req.SelectorID))

			if err != nil {
				c.JSON(500, dto.NewErrorAddSelectorResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessAddSelectorResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorAddSelectorResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorAddSelectorResponse(err))
		})
}

func (h *Handler) RemoveSelector(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorRemoveSelectorResponse(err))
		return
	}

	var req dto.RemoveSelectorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorRemoveSelectorResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorRemoveSelectorResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.RemoveSelector(id, uuid.MustParse(req.SelectorID))

			if err != nil {
				c.JSON(500, dto.NewErrorRemoveSelectorResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessRemoveSelectorResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorRemoveSelectorResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorRemoveSelectorResponse(err))
		})
}

func (h *Handler) AddFilter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorAddFilterResponse(err))
		return
	}

	var req dto.AddFilterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorAddFilterResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorAddFilterResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.AddFilter(id, uuid.MustParse(req.FilterID))

			if err != nil {
				c.JSON(500, dto.NewErrorAddFilterResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessAddFilterResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorAddFilterResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorAddFilterResponse(err))
		})
}

func (h *Handler) RemoveFilter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorRemoveFilterResponse(err))
		return
	}

	var req dto.RemoveFilterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorRemoveFilterResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorRemoveFilterResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.RemoveFilter(id, uuid.MustParse(req.FilterID))

			if err != nil {
				c.JSON(500, dto.NewErrorRemoveFilterResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessRemoveFilterResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorRemoveFilterResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorRemoveFilterResponse(err))
		})
}

func (h *Handler) AddField(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorAddFieldResponse(err))
		return
	}

	var req dto.AddFieldRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorAddFieldResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorAddFieldResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.AddField(id, uuid.MustParse(req.FieldID))

			if err != nil {
				c.JSON(500, dto.NewErrorAddFieldResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessAddFieldResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorAddFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorAddFieldResponse(err))
		})
}

func (h *Handler) RemoveField(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorRemoveFieldResponse(err))
		return
	}

	var req dto.RemoveFieldRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorRemoveFieldResponse(err))
		return
	}

	strat, err := h.service.Get(id)

	if err != nil {
		c.JSON(404, dto.NewErrorRemoveFieldResponse(err))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), strat).
		Allow(func() {

			err = h.service.RemoveField(id, uuid.MustParse(req.FieldID))

			if err != nil {
				c.JSON(500, dto.NewErrorRemoveFieldResponse(err))
				return
			}

			c.JSON(204, dto.NewSuccessRemoveFieldResponse())
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorRemoveFieldResponse(errors.New(reason)))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorRemoveFieldResponse(err))
		})
}
