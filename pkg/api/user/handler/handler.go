package handler

import (
	"juno/pkg/api/user"
	"juno/pkg/api/user/dto"
	"juno/pkg/api/user/policy"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger  logrus.FieldLogger
	userSvc user.Service
}

func New(logger logrus.FieldLogger, userSvc user.Service) *Handler {
	return &Handler{
		logger:  logger,
		userSvc: userSvc,
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

	found, err := h.userSvc.Get(id)

	if err != nil {
		h.logger.Error(err)
		c.JSON(404, dto.NewErrorGetUserResponse(
			user.ErrNotFound.Error(),
		))
		return
	}

	policy.CanRead(c.Request.Context(), found).
		Yes(func() {
			c.JSON(200, dto.NewSuccessGetUserResponse(found))
		}).
		No(func(reason string) {
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
