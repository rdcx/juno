package handler

import (
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/crawl/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	crawlService crawl.Service
}

func New(
	crawlService crawl.Service,
) *Handler {
	return &Handler{crawlService: crawlService}
}

func (h *Handler) Crawl(c *gin.Context) {
	var req dto.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
