package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/node"
	"juno/pkg/api/node/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger      *logrus.Logger
	policy      node.Policy
	nodeService node.Service
}

func New(l *logrus.Logger, policy node.Policy, ns node.Service) *Handler {
	return &Handler{
		logger:      l,
		policy:      policy,
		nodeService: ns,
	}
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	nodes, err := h.nodeService.ListByOwnerID(u.ID)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"user":  u.ID,
			})
		c.JSON(404, dto.NewErrorListNodesResponse(
			node.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanList(c.Request.Context(), nodes).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListNodesResponse(nodes))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error": reason,
					"user":  u.ID,
				})
			c.JSON(404, dto.NewErrorListNodesResponse(
				node.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorListNodesResponse(
				node.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	n, err := h.nodeService.Get(id)

	if err != nil {
		h.logger.Debug(
			logrus.Fields{
				"error": err.Error(),
				"node":  id,
			})
		c.JSON(404, dto.NewErrorGetNodeResponse(
			node.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanRead(c.Request.Context(), n).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetNodeResponse(n))
		}).
		Deny(func(reason string) {
			h.logger.Debug(
				logrus.Fields{
					"error": reason,
					"node":  id,
				})
			c.JSON(404, dto.NewErrorGetNodeResponse(
				node.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorGetNodeResponse(
				node.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Create(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateNodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			n, err := h.nodeService.Create(u.ID, req.Address, req.Shards)

			if err != nil {
				h.logger.Debug(
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
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorCreateNodeResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorCreateNodeResponse(
				node.ErrInternal.Error(),
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

	n, err := h.nodeService.Get(id)

	if err != nil {
		h.logger.Debug(
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

	h.policy.CanUpdate(c.Request.Context(), n).
		Allow(func() {

			var req dto.UpdateNodeRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err := req.ToDomain()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			n, err = h.nodeService.Update(id, n)

			if err != nil {
				h.logger.Debug(
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
			c.JSON(200, dto.NewSuccessUpdateNodeResponse(n))
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorUpdateNodeResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorUpdateNodeResponse(
				node.ErrInternal.Error(),
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

	n, err := h.nodeService.Get(id)

	if err != nil {
		h.logger.Debug(
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

	h.policy.CanDelete(c.Request.Context(), n).
		Allow(func() {
			err := h.nodeService.Delete(id)

			if err != nil {
				h.logger.Debug(
					logrus.Fields{
						"error": err.Error(),
						"user":  u.ID,
						"node":  id,
					})
				c.JSON(400, dto.NewErrorGetNodeResponse(
					err.Error(),
				))
				return
			}
			c.JSON(200, dto.NewSuccessDeleteNodeResponse())
		}).
		Deny(func(reason string) {
			c.JSON(401, dto.NewErrorDeleteNodeResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Debug(err)
			c.JSON(500, dto.NewErrorDeleteNodeResponse(
				node.ErrInternal.Error(),
			))
		})
}
