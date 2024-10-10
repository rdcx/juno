package handler

import (
	"juno/pkg/crawl/domain"
	"juno/pkg/shard/service/loadbalance"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoadBalanceHandler struct {
	lb *loadbalance.LoadBalancer
}

func NewLoadBalanceHandler(lb *loadbalance.LoadBalancer) *LoadBalanceHandler {
	return &LoadBalanceHandler{lb: lb}
}

func (h *LoadBalanceHandler) Crawl(c *gin.Context) {
	var req domain.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go h.lb.Crawl(req.URL)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
