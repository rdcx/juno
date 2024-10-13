package handler

import (
	"juno/pkg/api/user"
	"juno/pkg/api/user/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger      logrus.FieldLogger
	policy      user.Policy
	userService user.Service
}

func New(logger logrus.FieldLogger, policy user.Policy, userService user.Service) *Handler {
	return &Handler{
		logger:      logger,
		policy:      policy,
		userService: userService,
	}
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		h.logger.Error(err)
		c.JSON(400, dto.NewErrorGetUserResponse(
			user.ErrInvalidID.Error(),
		))
		return
	}

	found, err := h.userService.Get(id)

	if err != nil {
		h.logger.Error(err)
		c.JSON(404, dto.NewErrorGetUserResponse(
			user.ErrNotFound.Error(),
		))
		return
	}

	h.policy.CanRead(c.Request.Context(), found).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetUserResponse(found))
		}).
		Deny(func(reason string) {
			c.JSON(404, dto.NewErrorGetUserResponse(
				user.ErrNotFound.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Error(err)
			c.JSON(500, dto.NewErrorGetUserResponse(
				user.ErrInternal.Error(),
			))
		})
}
