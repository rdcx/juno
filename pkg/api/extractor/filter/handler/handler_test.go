package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/filter/dto"
	"juno/pkg/api/user"
	"juno/pkg/can"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockService struct {
	returnFilter  *filter.Filter
	returnFilters []*filter.Filter
	returnError   error
}

func (m mockService) Create(userID, fieldID uuid.UUID, name string, fType filter.FilterType, value string) (*filter.Filter, error) {
	return m.returnFilter, m.returnError
}

func (m mockService) Get(id uuid.UUID) (*filter.Filter, error) {
	return m.returnFilter, m.returnError
}

func (m mockService) ListByUserID(userID uuid.UUID) ([]*filter.Filter, error) {
	return m.returnFilters, m.returnError
}

type mockPolicy struct {
	allowed bool
	err     error
	reason  string
}

func (m mockPolicy) CanRead(ctx context.Context, sel *filter.Filter) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanCreate() can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanUpdate(ctx context.Context, sel *filter.Filter) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanDelete(ctx context.Context, sel *filter.Filter) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanList(ctx context.Context, sels []*filter.Filter) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &filter.Filter{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			FieldID: uuid.New(),
			Name:    "Value equals Charger",
			Type:    filter.FilterTypeStringEquals,
			Value:   "charger",
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnFilter: expect,
		})

		req := dto.CreateFilterRequest{
			Name:    expect.Name,
			FieldID: expect.FieldID.String(),
			Type:    string(expect.Type),
			Value:   expect.Value,
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/filters", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		handler.Create(c)

		if w.Code != 201 {
			t.Errorf("Expected 201, got %d", w.Code)
		}

		var resp dto.CreateFilterResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Filter == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Filter.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Filter.ID)
		}

		if resp.Filter.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Filter.Name)
		}

		if resp.Filter.Type != string(expect.Type) {
			t.Errorf("Expected %s, got %s", string(expect.Type), resp.Filter.Type)
		}

		if resp.Filter.Value != expect.Value {
			t.Errorf("Expected %s, got %s", expect.Value, resp.Filter.Value)
		}

	})

	t.Run("bad request", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/filters", strings.NewReader("bad json")).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 400 {
			t.Errorf("Expected 400, got %d", w.Code)
		}

		var resp dto.CreateFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{})

		req := dto.CreateFilterRequest{
			Name:    "Value equals Charger",
			FieldID: uuid.New().String(),
			Type:    string(filter.FilterTypeStringEquals),
			Value:   "charger",
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/filters", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.CreateFilterResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "reason" {
			t.Errorf("Expected reason, got %s", resp.Message)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &filter.Filter{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "Value equals Charger",
			Type:   filter.FilterTypeStringEquals,
			Value:  "charger",
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnFilter: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: expect.ID.String()}}

		c.Request = httptest.NewRequest("GET", "/filters/"+expect.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.GetFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Filter == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Filter.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Filter.ID)
		}

		if resp.Filter.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Filter.Name)
		}

		if resp.Filter.Type != string(expect.Type) {
			t.Errorf("Expected %s, got %s", string(expect.Type), resp.Filter.Type)
		}

		if resp.Filter.Value != expect.Value {
			t.Errorf("Expected %s, got %s", expect.Value, resp.Filter.Value)
		}
	})

	t.Run("not found", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnError: filter.ErrNotFound,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

		c.Request = httptest.NewRequest("GET", "/filters/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 404 {
			t.Errorf("Expected 404, got %d", w.Code)
		}

		var resp dto.GetFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != filter.ErrNotFound.Error() {
			t.Errorf("Expected %v, got %v", filter.ErrNotFound, resp.Message)
		}

	})

	t.Run("bad request", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "bad-uuid"}}

		c.Request = httptest.NewRequest("GET", "/filters/bad-uuid", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 400 {
			t.Errorf("Expected 400, got %d", w.Code)
		}

		var resp dto.GetFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		id := uuid.New()
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{
			returnFilter: &filter.Filter{
				ID:     id,
				UserID: uuid.New(),
			},
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		c.Request = httptest.NewRequest("GET", "/filters/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.GetFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "reason" {
			t.Errorf("Expected reason, got %s", resp.Message)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := []*filter.Filter{
			{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Name:   "Value equals Charger",
				Type:   filter.FilterTypeStringEquals,
				Value:  "charger",
			},
			{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Name:   "Value equals Charger",
				Type:   filter.FilterTypeStringEquals,
				Value:  "charger",
			},
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnFilters: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/filters", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.ListFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Filters == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if len(resp.Filters) != len(expect) {
			t.Errorf("Expected %d, got %d", len(expect), len(resp.Filters))
		}

		for i, f := range resp.Filters {
			if f.ID != expect[i].ID.String() {
				t.Errorf("Expected %s, got %s", expect[i].ID.String(), f.ID)
			}

			if f.Name != expect[i].Name {
				t.Errorf("Expected %s, got %s", expect[i].Name, f.Name)
			}

			if f.Type != string(expect[i].Type) {
				t.Errorf("Expected %s, got %s", string(expect[i].Type), f.Type)
			}

			if f.Value != expect[i].Value {
				t.Errorf("Expected %s, got %s", expect[i].Value, f.Value)
			}
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/filters", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.ListFilterResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != "reason" {
			t.Errorf("Expected reason, got %s", resp.Message)
		}
	})
}
