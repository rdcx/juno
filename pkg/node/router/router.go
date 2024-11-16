package router

import (
	"juno/pkg/node/crawl"
	"juno/pkg/node/extraction"
	"juno/pkg/node/info"

	"github.com/gin-gonic/gin"
)

func New(
	crawlHandler crawl.Handler,
	extractionHandler extraction.Handler,
	infoHandler info.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", crawlHandler.Crawl)
	r.POST("/extract", extractionHandler.Extract)
	r.GET("/info", infoHandler.Info)

	return r
}
