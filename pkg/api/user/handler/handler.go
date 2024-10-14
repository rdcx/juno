package handler

import (
	"juno/pkg/api/auth"
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

func (h *Handler) Profile(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	h.policy.CanRead(c.Request.Context(), u).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetUserResponse(u))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetUserResponse(
				user.ErrForbidden.Error(),
			))
		}).
		Err(func(err error) {
			h.logger.Error(err)
			c.JSON(500, dto.NewErrorGetUserResponse(
				user.ErrInternal.Error(),
			))
		})
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, dto.NewErrorCreateUserResponse(
			err.Error(),
		))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			u, err := h.userService.Create(req.Email, req.Password)

			if err != nil {
				h.logger.Error(err)
				c.JSON(400, dto.NewErrorCreateUserResponse(
					err.Error(),
				))
				return
			}

			c.JSON(201, dto.NewSuccessCreateUserResponse(u))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateUserResponse(
				reason,
			))
		}).
		Err(func(err error) {
			h.logger.Error(err)
			c.JSON(500, dto.NewErrorCreateUserResponse(
				user.ErrInternal.Error(),
			))
		})
}
