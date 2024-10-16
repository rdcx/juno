package handler

import (
	"encoding/json"
	"fmt"
	"juno/pkg/balancer/crawl/dto"
	"juno/pkg/balancer/crawl/service"
	"juno/pkg/shard"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/sirupsen/logrus"
)

func TestCrawl(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/crawl").
			Times(1).
			Reply(200).
			JSON(gin.H{"message": "ok"})

		svc := service.New(service.WithLogger(logrus.New()))
		svc.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})
		h := New(svc)

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

		// allow for the crawl goroutine to finish
		time.Sleep(10 * time.Millisecond)

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
