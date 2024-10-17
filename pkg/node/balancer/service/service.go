package service

import (
	"errors"
	apiClient "juno/pkg/api/client"
	balancerClient "juno/pkg/balancer/client"
	"juno/pkg/shard"
	"juno/pkg/url"
	"sync"
	"time"

	"math/rand"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger        *logrus.Logger
	apiClient     *apiClient.Client
	balancers     [shard.SHARDS][]string
	balancersLock sync.Mutex
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

	s.balancersLock.Lock()
	defer s.balancersLock.Unlock()

	for shardNum := range s.balancers {
		s.balancers[shardNum] = []string{}
	}

	for shardNum, b := range res.Shards {
		s.balancers[shardNum] = b
	}
}

func (s *Service) SetBalancers(balancers [shard.SHARDS][]string) {
	s.balancers = balancers
}

func (s *Service) ReportURLProcessed(urlStr string, status int) error {
	// TODO: implement

	return nil
}

func randomisedBalancersList(balancers []string) []string {

	rand.Shuffle(len(balancers), func(i, j int) {
		balancers[i], balancers[j] = balancers[j], balancers[i]
	})

	return balancers
}

func (s *Service) SendBatchedLinks(links []string) error {
	// group by shard
	groupedLinks := map[int][]string{}

	for _, link := range links {
		host, err := url.ToHostname(link)

		if err != nil {
			if s.logger != nil {
				s.logger.Error(err)
			}
			continue
		}

		shardNum := shard.GetShard(host)
		groupedLinks[shardNum] = append(groupedLinks[shardNum], link)
	}

	for shardNum, links := range groupedLinks {
		balancers := randomisedBalancersList(s.balancers[shardNum])

		if len(balancers) == 0 {
			if s.logger != nil {
				s.logger.Error("no balancers found for shard")
			}
			continue
		}

		for _, b := range balancers {
			balancerClient := balancerClient.New(
				"http://" + b,
			)

			err := balancerClient.CrawlURLs(links)

			if err != nil {
				if s.logger != nil {
					s.logger.Error(err)
					continue
				}
			}

			return nil
		}
	}

	return errors.New("failed to send batched links")
}

func (s *Service) SendCrawlRequest(urlStr string) error {
	host, err := url.ToHostname(urlStr)

	if err != nil {
		if s.logger != nil {
			s.logger.Error(err)
		}
		return err
	}

	shardNum := shard.GetShard(host)
	balancers := randomisedBalancersList(s.balancers[shardNum])

	if len(balancers) == 0 {
		if s.logger != nil {
			s.logger.Error("no balancers found for shard")
		}
		return errors.New("no balancers found for shard")
	}

	for _, b := range balancers {
		balancerClient := balancerClient.New(
			"http://" + b,
		)

		err := balancerClient.Crawl(urlStr)

		if err != nil {
			if s.logger != nil {
				s.logger.Error(err)
				continue
			}
		}

		return nil
	}

	return errors.New("failed to send crawl request")
}
