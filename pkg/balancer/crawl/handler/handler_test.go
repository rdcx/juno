package handler

import (
	"encoding/json"
	"fmt"
	"juno/pkg/balancer/crawl/dto"
	crawlService "juno/pkg/balancer/crawl/service"
	queueRepo "juno/pkg/balancer/queue/repo/mem"
	queueService "juno/pkg/balancer/queue/service"

	robotstxtRepo "juno/pkg/balancer/robotstxt/repo/mem"
	robotstxtService "juno/pkg/balancer/robotstxt/service"

	"juno/pkg/shard"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestCrawlURLs(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {

		repo := queueRepo.New()
		queueSvc := queueService.New(logrus.New(), repo)

		h := New(logrus.New(), queueSvc, robotstxtService.New(robotstxtRepo.New()))

		req := dto.CrawlURLsRequest{
			URLs: []string{"http://example.com"},
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/crawl/urls", strings.NewReader(fmt.Sprintf(`{"urls": ["%s"]}`, req.URLs[0])))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		h.CrawlURLs(c)

		// Then
		if c.Writer.Status() != http.StatusOK {
			t.Errorf("expected status 200 but got %d", c.Writer.Status())
		}

		var res dto.CrawlResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if res.Status != dto.OK {
			t.Errorf("expected %s but got %s", dto.OK, res.Status)
		}

		// check url has been added to queue
		pop, err := repo.Pop()
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if pop != req.URLs[0] {
			t.Errorf("expected %s but got %s", req.URLs[0], pop)
		}
	})
}

func TestCrawl(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {

		repo := queueRepo.New()
		queueSvc := queueService.New(logrus.New(), repo)
		svc := crawlService.New(crawlService.WithLogger(logrus.New()), crawlService.WithQueueService(queueSvc))
		svc.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})
		h := New(logrus.New(), queueSvc, robotstxtService.New(robotstxtRepo.New()))

		req := dto.CrawlRequest{
			URL: "http://example.com",
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/crawl", strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, req.URL)))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		h.Crawl(c)

		// Then
		if c.Writer.Status() != http.StatusOK {
			t.Errorf("expected status 200 but got %d", c.Writer.Status())
		}

		var res dto.CrawlResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if res.Status != dto.OK {
			t.Errorf("expected %s but got %s", dto.OK, res.Status)
		}

		// check url has been added to queue
		pop, err := repo.Pop()
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if pop != req.URL {
			t.Errorf("expected %s but got %s", req.URL, pop)
		}
	})
}
