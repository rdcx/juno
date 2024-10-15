package handler

import (
	"juno/pkg/api/assignment"
	"juno/pkg/api/assignment/dto"
	"juno/pkg/api/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger  *logrus.Logger
	policy  assignment.Policy
	service assignment.Service
}

func New(logger *logrus.Logger, policy assignment.Policy, service assignment.Service) *Handler {
	return &Handler{
		logger:  logger,
		policy:  policy,
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {

	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateAssignmentResponse(err))
		return
	}

	parsedNodeID, err := uuid.Parse(req.NodeID)

	if err != nil {
		c.JSON(400, dto.NewErrorCreateAssignmentResponse(err))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			assignment, err := h.service.Create(u.ID, parsedNodeID, req.Offset, req.Length)

			if err != nil {
				h.logger.WithError(err).Error("failed to create assignment")
				c.JSON(400, dto.NewErrorCreateAssignmentResponse(err))
				return
			}

			c.JSON(200, dto.NewSuccessCreateAssignmentResponse(assignment))
		}).
		Deny(func(reason string) {
			h.logger.WithField("reason", reason).Error("failed to create assignment")
			c.JSON(403, dto.NewErrorCreateAssignmentResponse(assignment.ErrForbidden))
		}).
		Err(func(err error) {
			h.logger.WithError(err).Error("failed to create assignment")
			c.JSON(500, dto.NewErrorCreateAssignmentResponse(assignment.ErrInternal))
		})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorGetAssignmentResponse(err))
		return
	}

	a, err := h.service.Get(id)

	if err != nil {
		h.logger.WithError(err).Error("failed to get assignment")
		c.JSON(404, dto.NewErrorGetAssignmentResponse(err))
		return
	}

	h.policy.CanRead(c.Request.Context(), a).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetAssignmentResponse(a))
		}).
		Deny(func(reason string) {
			h.logger.WithField("reason", reason).Error("failed to get assignment")
			c.JSON(403, dto.NewErrorGetAssignmentResponse(assignment.ErrForbidden))
		}).
		Err(func(err error) {
			h.logger.WithError(err).Error("failed to get assignment")
			c.JSON(500, dto.NewErrorGetAssignmentResponse(assignment.ErrInternal))
		})
}

func (h *Handler) ListByNodeID(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("node_id"))

	if err != nil {
		c.JSON(400, dto.NewErrorListAssignmentsResponse(err))
		return
	}

	assignments, err := h.service.ListByNodeID(nodeID)

	if err != nil {
		h.logger.WithError(err).Error("failed to list assignments")
		c.JSON(404, dto.NewErrorListAssignmentsResponse(assignment.ErrNotFound))
		return
	}

	h.policy.CanList(c.Request.Context(), assignments).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListAssignmentsResponse(assignments))
		}).
		Deny(func(reason string) {
			h.logger.WithField("reason", reason).Error("failed to list assignments")
			c.JSON(403, dto.NewErrorListAssignmentsResponse(assignment.ErrForbidden))
		}).
		Err(func(err error) {
			h.logger.WithError(err).Error("failed to list assignments")
			c.JSON(500, dto.NewErrorListAssignmentsResponse(assignment.ErrInternal))
		})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorUpdateAssignmentResponse(err))
		return
	}

	var req dto.UpdateAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorUpdateAssignmentResponse(err))
		return
	}

	a, err := h.service.Get(id)

	if err != nil {
		h.logger.WithError(err).Error("failed to update assignment")

		c.JSON(404, dto.NewErrorUpdateAssignmentResponse(assignment.ErrNotFound))
		return
	}

	h.policy.CanUpdate(c.Request.Context(), a).
		Allow(func() {
			assignment, err := h.service.Update(id, req.Offset, req.Length)

			if err != nil {
				h.logger.WithError(err).Error("failed to update assignment")
				c.JSON(400, dto.NewErrorUpdateAssignmentResponse(err))
				return
			}

			c.JSON(200, dto.NewSuccessUpdateAssignmentResponse(assignment))
		}).
		Deny(func(reason string) {
			h.logger.WithField("reason", reason).Error("failed to update assignment")
			c.JSON(403, dto.NewErrorUpdateAssignmentResponse(assignment.ErrForbidden))
		}).
		Err(func(err error) {
			h.logger.WithError(err).Error("failed to update assignment")
			c.JSON(500, dto.NewErrorUpdateAssignmentResponse(assignment.ErrInternal))
		})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, dto.NewErrorDeleteAssignmentResponse(err))
		return
	}

	a, err := h.service.Get(id)

	if err != nil {
		h.logger.WithError(err).Error("failed to delete assignment")
		c.JSON(404, dto.NewErrorDeleteAssignmentResponse(assignment.ErrNotFound))
		return
	}

	h.policy.CanDelete(c.Request.Context(), a).
		Allow(func() {
			if err := h.service.Delete(id); err != nil {
				h.logger.WithError(err).Error("failed to delete assignment")
				c.JSON(400, dto.NewErrorDeleteAssignmentResponse(err))
				return
			}

			c.JSON(200, dto.NewSuccessDeleteAssignmentResponse(a))
		}).
		Deny(func(reason string) {
			h.logger.WithField("reason", reason).Error("failed to delete assignment")
			c.JSON(403, dto.NewErrorDeleteAssignmentResponse(assignment.ErrForbidden))
		}).
		Err(func(err error) {
			h.logger.WithError(err).Error("failed to delete assignment")
			c.JSON(500, dto.NewErrorDeleteAssignmentResponse(assignment.ErrInternal))
		})
}
