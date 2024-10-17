package handler

import (
	"bytes"
	"encoding/json"
	"juno/pkg/node/runner/dto"
	"juno/pkg/node/runner/service"
	"net/http"
	"net/http/httptest"
	"testing"

	monkeyService "juno/pkg/monkey/service"
	pageRepo "juno/pkg/node/page/repo/mem"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestExecute(t *testing.T) {

	logger := logrus.New()
	pageRepo := pageRepo.New()
	pageService := pageService.New(pageRepo)
	storageService := storageService.New(
		t.TempDir(),
	)
	monkeyService := monkeyService.New()

	service := service.New(
		logger,
		pageService,
		storageService,
		monkeyService,
	)

	handler := New(logger, service)

	t.Run("success", func(t *testing.T) {
		src := `[0, "hello world"]`
		req := dto.ExecuteRequest{
			Src: src,
		}

		j, err := json.Marshal(req)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(j))
		handler.Execute(c)

		var res dto.ExecuteResponse

		err = json.NewDecoder(w.Body).Decode(&res)

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if string(res.Data) != `[0, "hello world"]` {
			t.Errorf("expected Hello, World! but got %s", string(res.Data))
		}
	})
}
