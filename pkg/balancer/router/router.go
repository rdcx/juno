package router

import (
	"juno/pkg/crawlbalance/balancer"

	"github.com/gin-gonic/gin"
)

func New(
	crawlHandler balancer.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", crawlHandler.Crawl)

	return r
}
