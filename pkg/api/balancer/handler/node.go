package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/balancer"
	"juno/pkg/api/balancer/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger          *logrus.Logger
	policy          balancer.Policy
	balancerService balancer.Service
}

func New(l *logrus.Logger, policy balancer.Policy, ns balancer.Service) *Handler {
	return &Handler{
		logger:          l,
		policy:          policy,
		balancerService: ns,
	}
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	balancers, err := h.balancerService.ListByOwnerID(u.ID)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
			})
		c.JSON(404, dto.NewErrorListBalancersResponse(
			balancer.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanList(c.Request.Context(), balancers).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListBalancersResponse(balancers))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error": reason,
					"user":  u.ID,
				})
			c.JSON(404, dto.NewErrorListBalancersResponse(
				balancer.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorListBalancersResponse(
				balancer.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.balancerService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error":    err.Error(),
				"balancer": id,
			})
		c.JSON(404, dto.NewErrorGetBalancerResponse(
			balancer.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanRead(c.Request.Context(), n).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetBalancerResponse(n))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error":    reason,
					"balancer": id,
				})
			c.JSON(404, dto.NewErrorGetBalancerResponse(
				balancer.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorGetBalancerResponse(
				balancer.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Create(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateBalancerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			n, err := h.balancerService.Create(u.ID, req.Address, req.Shards)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"req":   req,
					})
				c.JSON(400, dto.NewErrorGetBalancerResponse(
					err.Error(),
				))
				return
			}
			c.JSON(201, dto.NewSuccessGetBalancerResponse(n))
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorCreateBalancerResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorCreateBalancerResponse(
				balancer.ErrInternal.Error(),
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

	n, err := h.balancerService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error":    err.Error(),
				"user":     u.ID,
				"balancer": id,
			})
		c.JSON(404, dto.NewErrorGetBalancerResponse(
			balancer.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), n).
		Allow(func() {

			var req dto.UpdateBalancerRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err := req.ToDomain()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err = h.balancerService.Update(id, n)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"req":   req,
					})
				c.JSON(400, dto.NewErrorGetBalancerResponse(
					err.Error(),
				))
				return
			}
			c.JSON(200, dto.NewSuccessUpdateBalancerResponse(n))
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorUpdateBalancerResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorUpdateBalancerResponse(
				balancer.ErrInternal.Error(),
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

	n, err := h.balancerService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error":    err.Error(),
				"user":     u.ID,
				"balancer": id,
			})
		c.JSON(404, dto.NewErrorGetBalancerResponse(
			balancer.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanDelete(c.Request.Context(), n).
		Allow(func() {
			err := h.balancerService.Delete(id)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error":    err.Error(),
						"user":     u.ID,
						"balancer": id,
					})
				c.JSON(400, dto.NewErrorGetBalancerResponse(
					err.Error(),
				))
				return
			}
			c.JSON(200, dto.NewSuccessDeleteBalancerResponse())
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorDeleteBalancerResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorDeleteBalancerResponse(
				balancer.ErrInternal.Error(),
			))
		})
}
