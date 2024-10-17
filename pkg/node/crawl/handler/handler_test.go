package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/sirupsen/logrus"
)

type mockCrawlService struct {
	withError error
}

func (m *mockCrawlService) Crawl(ctx context.Context, url string) error {
	return m.withError
}

func TestCrawl(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {

		svc := New(logrus.New(), &mockCrawlService{})

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest("POST", "/crawl", strings.NewReader(`{"url": "http://example.com"}`))

		svc.Crawl(tc)

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}

		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		if w.Body.String() != `{"status":"success"}` {
			t.Errorf(`Expected response body to be {"status":"success"}, got %s`, w.Body.String())
		}

	})

	t.Run("should return error", func(t *testing.T) {

		svc := New(logrus.New(), &mockCrawlService{withError: errors.New("mock error")})

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest("POST", "/crawl", strings.NewReader(`{"url": "http://example.com"}`))

		svc.Crawl(tc)

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}

		if w.Code != 400 {
			t.Errorf("Expected status code 500, got %d", w.Code)
		}

		if w.Body.String() != `{"status":"error","message":"mock error"}` {
			t.Errorf(`Expected response body to be {"status":"error","message":"mock error"}, got %s`, w.Body.String())
		}
	})
}
