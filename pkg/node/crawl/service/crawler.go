package service

import (
	"context"
	"juno/pkg/crawl/domain"
	"juno/pkg/link"
	"log"
)

const CRAWL_TIMEOUT = 10

type Crawler struct {
	Queue         *Queue
	PolicyManager *PolicyManager
}

func NewCrawler() *Crawler {
	return &Crawler{
		Queue:         NewQueue(),
		PolicyManager: NewPolicyManager(),
	}
}

func (m *Crawler) Work() {
	for {
		url := m.Queue.Pop()
		if url == "" {
			break
		}

		host, err := link.ToHostname(url)

		if err != nil {
			log.Printf("failed to parse hostname in Crawler.Work: %v", err)
			continue
		}

		if m.PolicyManager.CanCrawl(host) {
			ctx, cancel := context.WithTimeout(context.Background(), CRAWL_TIMEOUT)

			status, _, _, err := FetchPage(ctx, url)
			if err != nil {
				log.Printf("failed to fetch page in Crawler.Work: %v", err)
				cancel()
			}

			if status != 429 {
				continue
			}
		}
		m.Queue.Push(url)
	}
}

func (m *Crawler) HandleCrawlRequest(req domain.CrawlRequest) {
	m.Queue.Push(req.URL)
}
