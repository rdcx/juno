package crawl

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrNoNodesAvailableInShard = errors.New("no nodes available in shard")
var ErrTooManyTries = errors.New("too many tries")

type Handler interface {
	Crawl(c *gin.Context)
}

type Service interface {
	Crawl(url string) error
}
