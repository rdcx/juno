package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/ranag"
	"juno/pkg/api/ranag/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger       *logrus.Logger
	policy       ranag.Policy
	ranagService ranag.Service
}

func New(l *logrus.Logger, policy ranag.Policy, ns ranag.Service) *Handler {
	return &Handler{
		logger:       l,
		policy:       policy,
		ranagService: ns,
	}
}

func (h *Handler) AllShardsRanags(c *gin.Context) {
	shards, err := h.ranagService.AllShardsRanags()

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
			})
		c.JSON(404, dto.NewErrorAllShardsRanagsResponse(
			ranag.ErrNotFound.Error(),
		))
		return
	}

	c.JSON(200, dto.NewSuccessAllShardsRanagsResponse(shards))
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	ranags, err := h.ranagService.ListByOwnerID(u.ID)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
			})
		c.JSON(404, dto.NewErrorListRanagsResponse(
			ranag.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanList(c.Request.Context(), ranags).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListRanagsResponse(ranags))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error": reason,
					"user":  u.ID,
				})
			c.JSON(404, dto.NewErrorListRanagsResponse(
				ranag.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorListRanagsResponse(
				ranag.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.ranagService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"ranag": id,
			})
		c.JSON(404, dto.NewErrorGetRanagResponse(
			ranag.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanRead(c.Request.Context(), n).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetRanagResponse(n))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error": reason,
					"ranag": id,
				})
			c.JSON(404, dto.NewErrorGetRanagResponse(
				ranag.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorGetRanagResponse(
				ranag.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Create(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateRanagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			n, err := h.ranagService.Create(u.ID, req.Address, req.ShardAssignments)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"req":   req,
					})
				c.JSON(400, dto.NewErrorGetRanagResponse(
					err.Error(),
				))
				return
			}
			c.JSON(201, dto.NewSuccessGetRanagResponse(n))
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorCreateRanagResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorCreateRanagResponse(
				ranag.ErrInternal.Error(),
			))
		})

}

func (h *Handler) Update(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.ranagService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
				"ranag": id,
			})
		c.JSON(404, dto.NewErrorGetRanagResponse(
			ranag.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), n).
		Allow(func() {

			var req dto.UpdateRanagRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err := req.ToDomain()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err = h.ranagService.Update(id, n)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"req":   req,
					})
				c.JSON(400, dto.NewErrorGetRanagResponse(
					err.Error(),
				))
				return
			}
			c.JSON(200, dto.NewSuccessUpdateRanagResponse(n))
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorUpdateRanagResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorUpdateRanagResponse(
				ranag.ErrInternal.Error(),
			))
		})

}

func (h *Handler) Delete(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.ranagService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
				"ranag": id,
			})
		c.JSON(404, dto.NewErrorGetRanagResponse(
			ranag.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanDelete(c.Request.Context(), n).
		Allow(func() {
			err := h.ranagService.Delete(id)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"ranag": id,
					})
				c.JSON(400, dto.NewErrorGetRanagResponse(
					err.Error(),
				))
				return
			}
			c.JSON(200, dto.NewSuccessDeleteRanagResponse())
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorDeleteRanagResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorDeleteRanagResponse(
				ranag.ErrInternal.Error(),
			))
		})
}
