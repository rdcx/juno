package service

import (
	"context"
	apiClient "juno/pkg/api/client"
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/policy"
	"juno/pkg/balancer/queue"
	"juno/pkg/node/client"
	"juno/pkg/shard"
	"juno/pkg/url"
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

func WithQueueService(queueService queue.Service) func(s *Service) {
	return func(s *Service) {
		s.queueService = queueService
	}
}

func WithPolicyService(policyService policy.Service) func(s *Service) {
	return func(s *Service) {
		s.policyService = policyService
	}
}

type Service struct {
	logger        *logrus.Logger
	apiClient     *apiClient.Client
	shards        [shard.SHARDS][]string
	shardsLock    sync.Mutex
	queueService  queue.Service
	policyService policy.Service
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

func (s *Service) ProcessQueue(ctx context.Context) error {

	if s.queueService == nil {
		panic("queue service is required")
	}

	if s.policyService == nil {
		panic("policy service is required")
	}

	for {
		select {
		case <-ctx.Done():
			return queue.ErrProcessQueueCancelled
		default:
			u, err := s.queueService.Pop()

			if err == queue.ErrNoURLsInQueue {
				select {
				case <-ctx.Done():
					return queue.ErrProcessQueueCancelled
				case <-time.After(500 * time.Millisecond):
					// continue after sleep
				}
				continue
			} else if err != nil {
				s.logger.Errorf("failed to pop url from queue: %v", err)
				continue
			}

			hostname, err := url.ToHostname(u)

			if err != nil {
				s.logger.Errorf("failed to get hostname from url: %v", err)
				continue
			}

			pol, err := s.policyService.Get(hostname)

			if err == policy.ErrPolicyNotFound {
				pol = policy.New(hostname)
			} else if err != nil {
				s.logger.Errorf("failed to get policy for url: %v", err)
				continue
			}

			if !s.policyService.CanCrawl(pol) {
				err = s.queueService.Push(u)

				if err != nil {
					s.logger.Errorf("failed to push url to queue: %v", err)
				}

				continue
			}

			crawlErr := s.Crawl(u)

			if crawlErr != nil {
				s.logger.Errorf("failed to crawl url: %v", crawlErr)
			}
			err = s.policyService.RecordCrawl(pol)
			if err != nil {
				s.logger.Errorf("failed to set policy for url: %v", err)
			}
		}
	}
}

func (s *Service) Crawl(u string) error {
	hostname, err := url.ToHostname(u)
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
		err = client.SendCrawlRequest(node, u)
		if err == nil {
			return nil
		}

		tries++
	}

	s.logger.Errorf("failed to send link %s to shard: %v", u, err)

	return crawl.ErrTooManyTries
}
