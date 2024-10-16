package client

import (
	"encoding/json"
	"juno/pkg/api/node/dto"
	"net/http"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) GetShards() (*dto.AllShardsNodesResponse, error) {
	var res dto.AllShardsNodesResponse
	resp, err := http.Get(c.baseURL + "/shards")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
