package service

import (
	apiClient "juno/pkg/api/client"
	nodeClient "juno/pkg/node/client"

	"juno/pkg/balancer/crawl"
	nodeDto "juno/pkg/node/dto"
	"juno/pkg/ranag/dto"
	"juno/pkg/shard"
	"math/rand"
	"sync"
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
	logger     *logrus.Logger
	apiClient  *apiClient.Client
	shards     [shard.SHARDS][]string
	shardsLock sync.Mutex
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

	s.shardsLock.Lock()
	defer s.shardsLock.Unlock()

	for i := 0; i < shard.SHARDS; i++ {
		s.shards[i] = []string{}
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
	s.shardsLock.Lock()
	defer s.shardsLock.Unlock()
	s.shards = shards
}

func (s *Service) randomNode(shard int) (string, error) {
	s.shardsLock.Lock()
	defer s.shardsLock.Unlock()
	if len(s.shards[shard]) == 0 {
		return "", crawl.ErrNoNodesAvailableInShard
	}

	return s.shards[shard][rand.Intn(len(s.shards[shard]))], nil
}

func (s *Service) QueryRange(offset int, total int, req dto.RangeAggregatorRequest) ([]map[string]interface{}, error) {

	data := make([]map[string]interface{}, 0)

	// get shard each shard for the given offset
	for i := offset; i < offset+total; i++ {
		node, err := s.randomNode(i % shard.SHARDS)
		if err != nil {
			s.logger.Errorf("failed to get random node: %v", err)
			return nil, err
		}

		selectors := []*nodeDto.Selector{}
		fields := []*nodeDto.Field{}

		for _, s := range req.Selectors {
			selectors = append(selectors, &nodeDto.Selector{
				ID:    s.ID,
				Value: s.Value,
			})
		}

		for _, f := range req.Fields {
			fields = append(fields, &nodeDto.Field{
				SelectorID: f.SelectorID,
				Name:       f.Name,
			})
		}

		// send request to the node
		extractions, err := nodeClient.SendExtractionRequest(node, selectors, fields)

		if err != nil {
			s.logger.Errorf("failed to send request to node: %v", err)
			return nil, err
		}

		data = append(data, extractions...)
	}

	return data, nil
}
