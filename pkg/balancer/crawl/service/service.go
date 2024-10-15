package service

import (
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/crawl/dto"
	"juno/pkg/link"
	"juno/pkg/node/client"
	"juno/pkg/shard"
	"log"
	"math/rand"
)

type Service struct {
	nodes [shard.SHARDS][]string
}

func New() *Service {
	return &Service{}
}

func (lb *Service) SetNodes(nodes [shard.SHARDS][]string) {
	lb.nodes = nodes
}

func (lb *Service) randomNode(shard int) (string, error) {
	if len(lb.nodes[shard]) == 0 {
		return "", crawl.ErrNoNodesAvailableInShard
	}

	return lb.nodes[shard][rand.Intn(len(lb.nodes[shard]))], nil
}

func (lb *Service) Crawl(req dto.CrawlRequest) {

	hostname, err := link.ToHostname(req.URL)
	if err != nil {
		return
	}
	shard := shard.GetShard(hostname)

	tries := 0
	for tries < 3 {
		node, err := lb.randomNode(shard)
		if err == crawl.ErrNoNodesAvailableInShard {
			log.Printf("no nodes available in shard %d", shard)
			return
		}
		err = client.SendCrawlRequest(node, req.URL)
		if err == nil {
			break
		}

		tries++
	}

	log.Printf("failed to send link %s to shard: %v", req.URL, err)
}
