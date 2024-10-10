package loadbalance

import (
	"juno/pkg/crawl/client"
	"juno/pkg/link"
	"juno/pkg/shard/domain"
	"log"
	"math/rand"
)

type LoadBalancer struct {
	nodes [domain.SHARDS][]string
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (lb *LoadBalancer) SetNodes(nodes [domain.SHARDS][]string) {
	lb.nodes = nodes
}

func (lb *LoadBalancer) randomNode(shard int) (string, error) {
	if len(lb.nodes[shard]) == 0 {
		return "", domain.ErrNoNodesAvailableInShard
	}

	return lb.nodes[shard][rand.Intn(len(lb.nodes[shard]))], nil
}

func (lb *LoadBalancer) Crawl(url string) {

	hostname, err := link.ToHostname(url)
	if err != nil {
		return
	}
	shard := domain.GetShard(hostname)

	tries := 0
	for tries < 3 {
		node, err := lb.randomNode(shard)
		if err == domain.ErrNoNodesAvailableInShard {
			log.Printf("no nodes available in shard %d", shard)
			return
		}
		err = client.SendCrawlRequest(node, url)
		if err == nil {
			break
		}

		tries++
	}

	log.Printf("failed to send link %s to shard: %v", url, err)
}
