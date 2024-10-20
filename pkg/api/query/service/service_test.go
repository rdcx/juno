package service

import (
	"juno/pkg/api/query"
	"juno/pkg/api/query/repo/mem"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)

		basicQuery := &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "title",
				MatchType: query.ExactStringMatch,
			},
		}

		userID := uuid.New()

		q, err := service.Create(userID, query.BasicQueryType, basicQuery)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if q.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, q.UserID)
		}

		if q.QueryType != query.BasicQueryType {
			t.Errorf("Expected %s, got %s", query.BasicQueryType, q.QueryType)
		}

		if q.BasicQuery.Title.Value != basicQuery.Title.Value {
			t.Errorf("Expected %s, got %s", basicQuery.Title.Value, q.BasicQuery.Title.Value)
		}

		if q.BasicQuery.Title.MatchType != basicQuery.Title.MatchType {
			t.Errorf("Expected %s, got %s", basicQuery.Title.MatchType, q.BasicQuery.Title.MatchType)
		}

		if q.Status != query.PendingStatus {
			t.Errorf("Expected %s, got %s", query.PendingStatus, q.Status)
		}

		if q.BasicQueryVersion != "v1" {
			t.Errorf("Expected v1, got %s", q.BasicQueryVersion)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)

		basicQuery := &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "title",
				MatchType: query.ExactStringMatch,
			},
		}

		userID := uuid.New()

		q, err := service.Create(userID, query.BasicQueryType, basicQuery)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := service.Get(q.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != q.ID {
			t.Errorf("Expected %s, got %s", q.ID, check.ID)
		}

		if check.UserID != q.UserID {
			t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
		}

		if check.Status != q.Status {
			t.Errorf("Expected %s, got %s", q.Status, check.Status)
		}

		if check.QueryType != q.QueryType {
			t.Errorf("Expected %s, got %s", q.QueryType, check.QueryType)
		}

		if check.BasicQueryVersion != q.BasicQueryVersion {
			t.Errorf("Expected %s, got %s", q.BasicQueryVersion, check.BasicQueryVersion)
		}

		if check.BasicQuery.Title.Value != q.BasicQuery.Title.Value {
			t.Errorf("Expected %s, got %s", q.BasicQuery.Title.Value, check.BasicQuery.Title.Value)
		}

		if check.BasicQuery.Title.MatchType != q.BasicQuery.Title.MatchType {
			t.Errorf("Expected %s, got %s", q.BasicQuery.Title.MatchType, check.BasicQuery.Title.MatchType)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)

		basicQuery := &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "title",
				MatchType: query.ExactStringMatch,
			},
		}

		userID := uuid.New()

		q, err := service.Create(userID, query.BasicQueryType, basicQuery)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		q.Status = query.CompletedStatus

		err = service.Update(q)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := service.Get(q.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.Status != query.CompletedStatus {
			t.Errorf("Expected %s, got %s", query.CompletedStatus, check.Status)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)

		basicQuery := &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "title",
				MatchType: query.ExactStringMatch,
			},
		}

		userID := uuid.New()

		q, err := service.Create(userID, query.BasicQueryType, basicQuery)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		queries, err := service.ListByUserID(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(queries) != 1 {
			t.Errorf("Expected 1, got %d", len(queries))
		}

		if queries[0].ID != q.ID {
			t.Errorf("Expected %s, got %s", q.ID, queries[0].ID)
		}
	})
}
