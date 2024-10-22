package handler

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/extraction/job"
	"juno/pkg/api/extraction/job/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	jobService job.Service
	policy     job.Policy
}

func New(jobService job.Service, policy job.Policy) *Handler {
	return &Handler{
		jobService: jobService,
		policy:     policy,
	}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/extraction-jobs", h.Create)
	r.GET("/extraction-jobs/:id", h.Get)
	r.GET("/extraction-jobs", h.List)
}

func (h *Handler) Create(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	var req dto.CreateJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCreateJobResponse(err.Error()))
		return
	}

	h.policy.CanCreate().
		Allow(func() {
			job, err := h.jobService.Create(u.ID, uuid.MustParse(req.ExtractorID))

			if err != nil {
				c.JSON(400, dto.NewErrorCreateJobResponse(err.Error()))
				return
			}

			c.JSON(201, dto.NewSuccessCreateJobResponse(job))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorCreateJobResponse(reason))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorCreateJobResponse(err.Error()))
		})
}

func (h *Handler) Get(c *gin.Context) {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.NewErrorGetJobResponse("invalid job ID"))
		return
	}

	job, err := h.jobService.Get(jobID)
	if err != nil {
		c.JSON(404, dto.NewErrorGetJobResponse("job not found"))
		return
	}

	h.policy.CanGet(c.Request.Context(), job).
		Allow(func() {
			c.JSON(200, dto.NewSuccessGetJobResponse(job))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorGetJobResponse(reason))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorGetJobResponse(err.Error()))
		})
}

func (h *Handler) List(c *gin.Context) {
	u := auth.MustUserFromContext(c.Request.Context())

	jobs, err := h.jobService.ListByUserID(u.ID)
	if err != nil {
		c.JSON(500, dto.NewErrorListJobsResponse("failed to fetch jobs"))
		return
	}

	h.policy.CanList(c.Request.Context(), jobs).
		Allow(func() {
			c.JSON(200, dto.NewSuccessListJobsResponse(jobs))
		}).
		Deny(func(reason string) {
			c.JSON(403, dto.NewErrorListJobsResponse(reason))
		}).
		Err(func(err error) {
			c.JSON(500, dto.NewErrorListJobsResponse(err.Error()))
		})
}
