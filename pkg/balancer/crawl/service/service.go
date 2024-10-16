package service

import (
	apiClient "juno/pkg/api/client"
	"juno/pkg/balancer/crawl"
	"juno/pkg/link"
	"juno/pkg/node/client"
	"juno/pkg/shard"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func WithLogger(logger *logrus.Logger) func(s *Service) {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithApiClient(apiClient *apiClient.Client) func(s *Service) {
	return func(s *Service) {
		s.apiClient = apiClient
	}
}

func WithShardFetchInterval(interval time.Duration) func(s *Service) {
	return func(s *Service) {

		if s.apiClient == nil {
			panic("api client is required")
		}

		go func() {
			s.fetchShards()
			for {
				time.Sleep(interval)
				s.fetchShards()
			}
		}()
	}
}

type Service struct {
	logger    *logrus.Logger
	apiClient *apiClient.Client
	shards    [shard.SHARDS][]string
}

func New(options ...func(s *Service)) *Service {
	s := &Service{}

	for _, option := range options {
		option(s)
	}

	if s.logger == nil {
		panic("logger is required")
	}

	return s
}

func (s *Service) fetchShards() {
	res, err := s.apiClient.GetShards()
	if err != nil {
		s.logger.Printf("failed to fetch shards: %v", err)
		return
	}

	for shard, nodes := range res.Shards {
		s.shards[shard] = nodes
	}

	if len(res.Shards) != shard.SHARDS {
		s.logger.Errorf("expected %d shards, got %d", shard.SHARDS, len(res.Shards))
		return
	}

	s.logger.Infof("shards fetched: %d", len(res.Shards))
}

func (s *Service) SetShards(shards [shard.SHARDS][]string) {
	s.shards = shards
}

func (s *Service) randomNode(shard int) (string, error) {
	if len(s.shards[shard]) == 0 {
		return "", crawl.ErrNoNodesAvailableInShard
	}

	return s.shards[shard][rand.Intn(len(s.shards[shard]))], nil
}

func (s *Service) Crawl(url string) error {
	hostname, err := link.ToHostname(url)
	if err != nil {
		s.logger.Errorf("failed to parse hostname: %v", err)
		return err
	}
	shard := shard.GetShard(hostname)

	tries := 0
	for tries < 3 {
		node, err := s.randomNode(shard)
		if err == crawl.ErrNoNodesAvailableInShard {
			s.logger.Errorf("no nodes available in shard %d", shard)
			return err
		}
		err = client.SendCrawlRequest(node, url)
		if err == nil {
			break
		}

		tries++
	}

	s.logger.Errorf("failed to send link %s to shard: %v", url, err)

	return crawl.ErrTooManyTries
}
