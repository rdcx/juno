package loadbalance

import "github.com/gin-gonic/gin"

func shardRouter(lb *ShardLoadBalancer) *gin.Engine {
	r := gin.Default()

	r.POST("/crawl", lb.Crawl)

	return r
}

func RunShardLoadBalancer(port string) {
	shardLb := NewShardLoadBalancer()

	shardRouter(shardLb).Run(":" + port)
}
