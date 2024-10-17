package router

import (
	"juno/pkg/node/crawl"
	"juno/pkg/node/runner"

	"github.com/gin-gonic/gin"
)

func New(
	crawlHandler crawl.Handler,
	runnerHandler runner.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", crawlHandler.Crawl)
	r.POST("/runner/execute", runnerHandler.Execute)
	r.POST("/runner/titles", runnerHandler.Titles)

	return r
}
