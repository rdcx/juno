package client

import (
	"juno/pkg/api/node/dto"
	"testing"

	"github.com/h2non/gock"
)

func TestGetShards(t *testing.T) {
	t.Run("should return all shards", func(t *testing.T) {
		// Given
		baseURL := "http://localhost:8080"
		client := New(baseURL)
		expected := &dto.AllShardsNodesResponse{
			Shards: map[int][]string{
				1: {"node1", "node2"},
				2: {"node1", "node2"},
			},
		}

		defer gock.Off()

		gock.New(baseURL).
			Get("/shards").
			Reply(200).
			JSON(expected)

		res, err := client.GetShards()

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if res.Shards[1][0] != "node1" {
			t.Errorf("expected node1 but got %s", res.Shards[1][0])
		}

		if res.Shards[2][0] != "node1" {
			t.Errorf("expected node1 but got %s", res.Shards[2][0])
		}
	})

	t.Run("should return error when request fails", func(t *testing.T) {
		// Given
		baseURL := "http://localhost:8080"
		client := New(baseURL)

		defer gock.Off()

		gock.New(baseURL).
			Get("/shards").
			Reply(500)

		_, err := client.GetShards()

		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})
}
