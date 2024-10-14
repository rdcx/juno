package handler

import (
	loadbalance "juno/pkg/crawlbalance/balancer/service"
	"juno/pkg/crawlbalance/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	lb *loadbalance.LoadBalancer
}

func NewLoadBalanceHandler(lb *loadbalance.LoadBalancer) *Handler {
	return &Handler{lb: lb}
}

func (h *Handler) Crawl(c *gin.Context) {
	var req dto.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go h.lb.Crawl(req.URL)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
