package handler

import (
	"juno/pkg/balancer/crawl/dto"
	"juno/pkg/balancer/queue"
	"juno/pkg/balancer/robotstxt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger           *logrus.Logger
	queueService     queue.Service
	robotsTxtService robotstxt.Service
}

func New(
	logger *logrus.Logger,
	queueService queue.Service,
	robotsTxtService robotstxt.Service,
) *Handler {
	return &Handler{
		logger:           logger,
		queueService:     queueService,
		robotsTxtService: robotsTxtService,
	}
}

func (h *Handler) CrawlURLs(c *gin.Context) {
	var req dto.CrawlURLsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, url := range req.URLs {

		if !h.robotsTxtService.CanCrawlURL(url) {
			continue
		}

		if err := h.queueService.Push(url); err != nil {
			h.logger.WithError(err).Error("failed to push url to queue")
		}
	}

	c.JSON(http.StatusOK, dto.NewOKCrawlResponse())
}

func (h *Handler) Crawl(c *gin.Context) {
	var req dto.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !h.robotsTxtService.CanCrawlURL(req.URL) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to crawl"})
		return
	}

	if err := h.queueService.Push(req.URL); err != nil {
		h.logger.WithError(err).Error("failed to push url to queue")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.NewOKCrawlResponse())
}
