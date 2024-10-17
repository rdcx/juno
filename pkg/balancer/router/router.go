package router

import (
	"juno/pkg/balancer/crawl"

	"github.com/gin-gonic/gin"
)

func New(
	crawlHandler crawl.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", crawlHandler.Crawl)
	r.POST("/crawl/urls", crawlHandler.CrawlURLs)

	return r
}
