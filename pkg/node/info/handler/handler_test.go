package handler

import (
	"encoding/json"
	"errors"
	"juno/pkg/node/info"
	"juno/pkg/node/info/dto"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockInfoService struct {
	withErr bool
}

func (m *mockInfoService) GetInfo() (*info.Info, error) {
	if m.withErr {
		return nil, errors.New("error")
	}

	return &info.Info{
		PageCount: 10,
	}, nil
}

func TestHandler_Info(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := New(&mockInfoService{false})
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		h.Info(c)

		if c.Writer.Status() != 200 {
			t.Fatalf("expected 200, got %d", c.Writer.Status())
		}

		var response dto.InfoResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.Info.PageCount != 10 {
			t.Fatalf("expected 10, got %d", response.Info.PageCount)
		}
	})

	t.Run("error", func(t *testing.T) {
		h := New(&mockInfoService{true})
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		h.Info(c)

		if c.Writer.Status() != 500 {
			t.Fatalf("expected 500, got %d", c.Writer.Status())
		}

		var response dto.InfoResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.Status != dto.ERROR {
			t.Fatalf("expected %s, got %s", dto.ERROR, response.Status)
		}
	})
}
