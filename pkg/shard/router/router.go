package router

import (
	"juno/pkg/shard/handler"

	"github.com/gin-gonic/gin"
)

func shardRouter(lbHandler *handler.LoadBalanceHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", lbHandler.Crawl)

	return r
}

func RunShardService(lbHandler *handler.LoadBalanceHandler, port string) {
	shardRouter(
		lbHandler,
	).Run(":" + port)
}
