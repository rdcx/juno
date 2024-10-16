package service

import (
	apiClient "juno/pkg/api/client"
	balancerClient "juno/pkg/balancer/client"
	"juno/pkg/shard"
	"juno/pkg/url"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger    *logrus.Logger
	apiClient *apiClient.Client
	balancers [shard.SHARDS][]string
}

func WithApiClient(api *apiClient.Client) func(s *Service) {
	return func(s *Service) {
		s.apiClient = api
	}
}

func WithLogger(logger *logrus.Logger) func(s *Service) {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithBalancerFetchInterval(interval time.Duration) func(s *Service) {
	return func(s *Service) {
		if s.apiClient == nil {
			panic("api client is required")
		}

		go func() {
			s.fetchBalancers()
			for {
				time.Sleep(interval)
				s.fetchBalancers()
			}
		}()
	}
}

func New(options ...func(s *Service)) *Service {
	s := &Service{}

	for _, o := range options {
		o(s)
	}

	return s
}

func (s *Service) fetchBalancers() {
	res, err := s.apiClient.GetBalancers()

	if err != nil {
		return
	}

	for shardNum := range s.balancers {
		s.balancers[shardNum] = []string{}
	}

	for shardNum, b := range res.Shards {
		s.balancers[shardNum] = b
	}
}

func (s *Service) ReportURLProcessed(urlStr string) {
	// TODO: implement
}

func (s *Service) SendCrawlRequest(urlStr string) {
	host, err := url.ToHostname(urlStr)

	if err != nil {
		if s.logger != nil {
			s.logger.Error(err)
		}
		return
	}

	shardNum := shard.GetShard(host)
	balancers := s.balancers[shardNum]

	if len(balancers) == 0 {
		if s.logger != nil {
			s.logger.Error("no balancers found for shard")
		}
		return
	}

	for _, b := range balancers {
		balancerClient := balancerClient.New(b)

		err := balancerClient.Crawl(urlStr)

		if err != nil {
			if s.logger != nil {
				s.logger.Error(err)
			}
		}
	}
}
