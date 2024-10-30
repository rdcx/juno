package service

import (
	apiClient "juno/pkg/api/client"
	nodeClient "juno/pkg/node/client"

	"juno/pkg/balancer/crawl"
	extractionDto "juno/pkg/node/extraction/dto"
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

func (s *Service) RangeAggregate(offset int, total int, req dto.RangeAggregatorRequest) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, 0)
	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		errChan  = make(chan error, 1) // To capture errors from goroutines
		dataChan = make(chan []map[string]interface{}, total)
		sem      = make(chan struct{}, 10) // Buffered channel to limit to 10 concurrent workers
	)

	// Launch workers for each shard
	for shard := offset; shard < offset+total; shard++ {
		wg.Add(1)
		go func(shard int) {
			defer wg.Done()
			sem <- struct{}{} // Block if there are already 10 workers

			node, err := s.randomNode(shard)
			if err != nil {
				select {
				case errChan <- err: // Send error if no error has been sent
				default:
				}
				<-sem // Release a spot in the semaphore
				return
			}

			selectors := make([]*extractionDto.Selector, len(req.Selectors))
			for i, s := range req.Selectors {
				selectors[i] = &extractionDto.Selector{
					ID:    s.ID,
					Value: s.Value,
				}
			}

			fields := make([]*extractionDto.Field, len(req.Fields))
			for i, f := range req.Fields {
				fields[i] = &extractionDto.Field{
					SelectorID: f.SelectorID,
					Name:       f.Name,
				}
			}

			// Send request to the node
			extractions, err := nodeClient.SendExtractionRequest(node, shard, selectors, fields)
			if err != nil {
				s.logger.Errorf("failed to send request to node: %v", err)
				select {
				case errChan <- err:
				default:
				}
				<-sem
				return
			}

			dataChan <- extractions // Send extractions to data channel
			<-sem                   // Release a spot in the semaphore
		}(shard)
	}

	// Close the data channel once all goroutines complete
	go func() {
		wg.Wait()
		close(dataChan)
		close(errChan)
	}()

	// Collect results from dataChan
	for extraction := range dataChan {
		mu.Lock()
		data = append(data, extraction...)
		mu.Unlock()
	}

	// Check if any error occurred
	if err, ok := <-errChan; ok {
		return nil, err
	}

	return data, nil
}
