package crawl

import (
	"errors"
	"juno/pkg/balancer/crawl/dto"

	"github.com/gin-gonic/gin"
)

var ErrNoNodesAvailableInShard = errors.New("no nodes available in shard")

type Handler interface {
	Crawl(c *gin.Context)
}

type Service interface {
	Crawl(req dto.CrawlRequest) (dto.CrawlResponse, error)
}
