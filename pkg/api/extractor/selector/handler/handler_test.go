package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/api/extractor/selector/dto"
	"juno/pkg/api/user"
	"juno/pkg/can"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockService struct {
	returnSelector  *selector.Selector
	returnSelectors []*selector.Selector
	returnError     error
}

func (m mockService) Create(userID uuid.UUID, name, value string, visibility selector.Visibility) (*selector.Selector, error) {
	return m.returnSelector, m.returnError
}

func (m mockService) Get(id uuid.UUID) (*selector.Selector, error) {
	return m.returnSelector, m.returnError
}

func (m mockService) ListByUserID(userID uuid.UUID) ([]*selector.Selector, error) {
	return m.returnSelectors, m.returnError
}

type mockPolicy struct {
	allowed bool
	err     error
	reason  string
}

func (m mockPolicy) CanRead(ctx context.Context, sel *selector.Selector) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanCreate() can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanUpdate(ctx context.Context, sel *selector.Selector) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanDelete(ctx context.Context, sel *selector.Selector) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanList(ctx context.Context, sels []*selector.Selector) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &selector.Selector{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			Name:       "name",
			Value:      "#productTitle",
			Visibility: selector.VisibilityPrivate,
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnSelector: expect,
		})
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		req := dto.CreateSelectorRequest{
			Name:       name,
			Value:      value,
			Visibility: string(visibility),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/selectors", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		handler.Create(c)

		if w.Code != 201 {
			t.Errorf("Expected 201, got %d", w.Code)
		}

		var resp dto.CreateSelectorResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Selector == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Selector.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Selector.ID)
		}

		if resp.Selector.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Selector.Name)
		}

		if resp.Selector.Value != expect.Value {
			t.Errorf("Expected %s, got %s", expect.Value, resp.Selector.Value)
		}

		if resp.Selector.Visibility != string(expect.Visibility) {
			t.Errorf("Expected %s, got %s", string(expect.Visibility), resp.Selector.Visibility)
		}
	})

	t.Run("bad request", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{})

		req := dto.CreateSelectorRequest{}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/selectors", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 400 {
			t.Errorf("Expected 400, got %d", w.Code)
		}

		var resp dto.CreateSelectorResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if !strings.Contains(resp.Message, "Field validation") {
			t.Errorf("Expected Field validation, got %s", resp.Message)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "forbidden"}, mockService{})
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		req := dto.CreateSelectorRequest{
			Name:       name,
			Value:      value,
			Visibility: string(visibility),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/selectors", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.CreateSelectorResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "forbidden" {
			t.Errorf("Expected forbidden, got %s", resp.Message)
		}
	})

	t.Run("error", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnError: errors.New("error"),
		})
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		req := dto.CreateSelectorRequest{
			Name:       name,
			Value:      value,
			Visibility: string(visibility),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/selectors", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 500 {
			t.Errorf("Expected 500, got %d", w.Code)
		}

		var resp dto.CreateSelectorResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "error" {
			t.Errorf("Expected error, got %s", resp.Message)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &selector.Selector{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			Name:       "name",
			Value:      "#productTitle",
			Visibility: selector.VisibilityPrivate,
		}
		handler := New(mockPolicy{allowed: true}, mockService{
			returnSelector: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: expect.ID.String()}}

		c.Request = httptest.NewRequest("GET", "/selectors/"+expect.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.GetSelectorResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Selector == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Selector.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Selector.ID)
		}

		if resp.Selector.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Selector.Name)
		}

		if resp.Selector.Value != expect.Value {
			t.Errorf("Expected %s, got %s", expect.Value, resp.Selector.Value)
		}

		if resp.Selector.Visibility != string(expect.Visibility) {
			t.Errorf("Expected %s, got %s", string(expect.Visibility), resp.Selector.Visibility)
		}
	})

	t.Run("not found", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnError: selector.ErrNotFound,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}

		c.Request = httptest.NewRequest("GET", "/selectors/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 404 {
			t.Errorf("Expected 404, got %d", w.Code)
		}

		var resp dto.GetSelectorResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "selector not found" {
			t.Errorf("Expected selector not found, got %s", resp.Message)
		}

	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "forbidden"}, mockService{
			returnSelector: &selector.Selector{UserID: uuid.New()},
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}

		c.Request = httptest.NewRequest("GET", "/selectors/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.GetSelectorResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "forbidden" {
			t.Errorf("Expected forbidden, got %s", resp.Message)
		}
	})

}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := []*selector.Selector{
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Name:       "name",
				Value:      "#productTitle",
				Visibility: selector.VisibilityPrivate,
			},
		}
		handler := New(mockPolicy{allowed: true}, mockService{
			returnSelectors: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/selectors", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.ListSelectorResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if len(resp.Selectors) != 1 {
			t.Errorf("Expected 1, got %d", len(resp.Selectors))
		}

		if resp.Selectors[0].ID != expect[0].ID.String() {
			t.Errorf("Expected %s, got %s", expect[0].ID.String(), resp.Selectors[0].ID)
		}

		if resp.Selectors[0].Name != expect[0].Name {
			t.Errorf("Expected %s, got %s", expect[0].Name, resp.Selectors[0].Name)
		}

		if resp.Selectors[0].Value != expect[0].Value {
			t.Errorf("Expected %s, got %s", expect[0].Value, resp.Selectors[0].Value)
		}

		if resp.Selectors[0].Visibility != string(expect[0].Visibility) {
			t.Errorf("Expected %s, got %s", string(expect[0].Visibility), resp.Selectors[0].Visibility)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "forbidden"}, mockService{
			returnSelectors: []*selector.Selector{&selector.Selector{UserID: uuid.New()}},
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/selectors", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.ListSelectorResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "forbidden" {
			t.Errorf("Expected forbidden, got %s", resp.Message)
		}
	})
}
