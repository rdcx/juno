package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/extractor/strategy/dto"
	"juno/pkg/api/user"
	"juno/pkg/can"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockService struct {
	returnStrategy  *strategy.Strategy
	returnStrategys []*strategy.Strategy
	returnError     error
}

func (m mockService) Create(userID, selectorID uuid.UUID, name string, fType strategy.StrategyType) (*strategy.Strategy, error) {
	return m.returnStrategy, m.returnError
}

func (m mockService) Get(id uuid.UUID) (*strategy.Strategy, error) {
	return m.returnStrategy, m.returnError
}

func (m mockService) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	return m.returnStrategys, m.returnError
}

func (m mockService) ListBySelectorID(selectorID uuid.UUID) ([]*strategy.Strategy, error) {
	return m.returnStrategys, m.returnError
}

type mockPolicy struct {
	allowed bool
	err     error
	reason  string
}

func (m mockPolicy) CanRead(ctx context.Context, sel *strategy.Strategy) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanCreate() can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanUpdate(ctx context.Context, sel *strategy.Strategy) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanDelete(ctx context.Context, sel *strategy.Strategy) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func (m mockPolicy) CanList(ctx context.Context, sels []*strategy.Strategy) can.Result {
	return can.Result{Allowed: m.allowed, Error: m.err, Reason: m.reason}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &strategy.Strategy{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			SelectorID: uuid.New(),
			Name:       "product_title",
			Type:       strategy.StrategyTypeString,
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnStrategy: expect,
		})

		req := dto.CreateStrategyRequest{
			Name:       expect.Name,
			Type:       string(expect.Type),
			SelectorID: expect.SelectorID.String(),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/strategys", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		handler.Create(c)

		if w.Code != 201 {
			t.Errorf("Expected 201, got %d", w.Code)
		}

		var resp dto.CreateStrategyResponse

		err = json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Strategy == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Strategy.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Strategy.ID)
		}

		if resp.Strategy.SelectorID != expect.SelectorID.String() {
			t.Errorf("Expected %s, got %s", expect.SelectorID.String(), resp.Strategy.SelectorID)
		}

		if resp.Strategy.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Strategy.Name)
		}

		if resp.Strategy.Type != string(expect.Type) {
			t.Errorf("Expected %s, got %s", string(expect.Type), resp.Strategy.Type)
		}
	})

	t.Run("bad request", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/strategys", strings.NewReader("bad json")).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 400 {
			t.Errorf("Expected 400, got %d", w.Code)
		}

		var resp dto.CreateStrategyResponse

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

		req := dto.CreateStrategyRequest{
			Name:       "Value equals Charger",
			SelectorID: uuid.New().String(),
			Type:       string(strategy.StrategyTypeInteger),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/strategys", bytes.NewBuffer(encoded)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Create(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.CreateStrategyResponse

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
		expect := &strategy.Strategy{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			SelectorID: uuid.New(),
			Name:       "product_title",
			Type:       strategy.StrategyTypeString,
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnStrategy: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: expect.ID.String()})

		c.Request = httptest.NewRequest("GET", "/strategys/"+expect.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.GetStrategyResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if resp.Strategy == nil {
			t.Errorf("Expected non-nil, got nil")
		}

		if resp.Strategy.ID != expect.ID.String() {
			t.Errorf("Expected %s, got %s", expect.ID.String(), resp.Strategy.ID)
		}

		if resp.Strategy.SelectorID != expect.SelectorID.String() {
			t.Errorf("Expected %s, got %s", expect.SelectorID.String(), resp.Strategy.SelectorID)
		}

		if resp.Strategy.Name != expect.Name {
			t.Errorf("Expected %s, got %s", expect.Name, resp.Strategy.Name)
		}

		if resp.Strategy.Type != string(expect.Type) {
			t.Errorf("Expected %s, got %s", string(expect.Type), resp.Strategy.Type)
		}
	})

	t.Run("not found", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnError: strategy.ErrNotFound,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.New().String()})

		c.Request = httptest.NewRequest("GET", "/strategys/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 404 {
			t.Errorf("Expected 404, got %d", w.Code)
		}

		var resp dto.GetStrategyResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "error" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if resp.Message != strategy.ErrNotFound.Error() {
			t.Errorf("Expected %s, got %s", strategy.ErrNotFound.Error(), resp.Message)
		}

	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.New().String()})

		c.Request = httptest.NewRequest("GET", "/strategys/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.GetStrategyResponse

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
		expect := []*strategy.Strategy{
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				SelectorID: uuid.New(),
				Name:       "product_title",
				Type:       strategy.StrategyTypeString,
			},
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				SelectorID: uuid.New(),
				Name:       "product_price",
				Type:       strategy.StrategyTypeInteger,
			},
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnStrategys: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/strategys", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.ListStrategyResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if len(resp.Strategys) != len(expect) {
			t.Errorf("Expected %d, got %d", len(expect), len(resp.Strategys))
		}

		for i, f := range resp.Strategys {
			if f.ID != expect[i].ID.String() {
				t.Errorf("Expected %s, got %s", expect[i].ID.String(), f.ID)
			}

			if f.SelectorID != expect[i].SelectorID.String() {
				t.Errorf("Expected %s, got %s", expect[i].SelectorID.String(), f.SelectorID)
			}

			if f.Name != expect[i].Name {
				t.Errorf("Expected %s, got %s", expect[i].Name, f.Name)
			}

			if f.Type != string(expect[i].Type) {
				t.Errorf("Expected %s, got %s", string(expect[i].Type), f.Type)
			}
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/strategys", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.List(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.ListStrategyResponse

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

func TestListBySelectorID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := []*strategy.Strategy{
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				SelectorID: uuid.New(),
				Name:       "product_title",
				Type:       strategy.StrategyTypeString,
			},
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				SelectorID: uuid.New(),
				Name:       "product_price",
				Type:       strategy.StrategyTypeInteger,
			},
		}
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnStrategys: expect,
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/strategys/selector/"+expect[0].SelectorID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.ListBySelectorID(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.ListStrategyResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected success, got %s", resp.Status)
		}

		if len(resp.Strategys) != len(expect) {
			t.Errorf("Expected %d, got %d", len(expect), len(resp.Strategys))
		}

		for i, f := range resp.Strategys {
			if f.ID != expect[i].ID.String() {
				t.Errorf("Expected %s, got %s", expect[i].ID.String(), f.ID)
			}

			if f.SelectorID != expect[i].SelectorID.String() {
				t.Errorf("Expected %s, got %s", expect[i].SelectorID.String(), f.SelectorID)
			}

			if f.Name != expect[i].Name {
				t.Errorf("Expected %s, got %s", expect[i].Name, f.Name)
			}

			if f.Type != string(expect[i].Type) {
				t.Errorf("Expected %s, got %s", string(expect[i].Type), f.Type)
			}

		}

	})

	t.Run("forbidden", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: false, reason: "reason"}, mockService{})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/strategys/selector/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.ListBySelectorID(c)

		if w.Code != 403 {
			t.Errorf("Expected 403, got %d", w.Code)
		}

		var resp dto.ListStrategyResponse

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

	t.Run("empty", func(t *testing.T) {
		handler := New(&mockPolicy{allowed: true}, mockService{
			returnStrategys: []*strategy.Strategy{},
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.New().String()})

		c.Request = httptest.NewRequest("GET", "/selectors/"+uuid.New().String()+"/strategys", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.ListBySelectorID(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp dto.ListStrategyResponse

		err := json.NewDecoder(w.Body).Decode(&resp)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected error, got %s", resp.Status)
		}

		if len(resp.Strategys) != 0 {
			t.Errorf("Expected 0, got %d", len(resp.Strategys))
		}
	})
}
