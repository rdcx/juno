package service

type testCrawlService struct {
	returnErr error
	hits      []string
}

func (s *testCrawlService) Crawl(url string) error {
	if s.returnErr != nil {
		return s.returnErr
	}

	s.hits = append(s.hits, url)
	return nil
}
