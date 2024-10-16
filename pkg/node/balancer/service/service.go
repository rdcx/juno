package service

import (
	"juno/pkg/api/client"
	"juno/pkg/shard"
	"time"
)

type Service struct {
	apiClient *client.Client
	balancers [shard.SHARDS][]string
}

func WithApiClient(api *client.Client) func(s *Service) {
	return func(s *Service) {
		s.apiClient = api
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
