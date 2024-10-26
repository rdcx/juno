package router

import (
	"juno/pkg/node/crawl"
	"juno/pkg/node/extraction"

	"github.com/gin-gonic/gin"
)

func New(
	crawlHandler crawl.Handler,
	extractionHandler extraction.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", crawlHandler.Crawl)
	r.POST("/extract", extractionHandler.Extract)

	return r
}
