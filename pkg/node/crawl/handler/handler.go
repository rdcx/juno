package handler

import (
	"context"
	"juno/pkg/node/crawl"
	"juno/pkg/node/crawl/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger       *logrus.Logger
	crawlService crawl.Service
}

func New(logger *logrus.Logger, crawlService crawl.Service) *Handler {
	return &Handler{
		logger:       logger,
		crawlService: crawlService,
	}
}

func (h *Handler) Crawl(c *gin.Context) {

	var req dto.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorCrawlResponse(err.Error()))
		return
	}

	err := h.crawlService.Crawl(context.Background(), req.URL)

	if err != nil {
		c.JSON(400, dto.NewErrorCrawlResponse(err.Error()))
		return
	}

	c.JSON(200, dto.NewSuccessCrawlResponse())
}
