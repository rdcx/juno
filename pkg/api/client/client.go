package client

import (
	"encoding/json"
	balancerDto "juno/pkg/api/balancer/dto"
	nodeDto "juno/pkg/api/node/dto"
	"net/http"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) GetShards() (*nodeDto.AllShardsNodesResponse, error) {
	var res nodeDto.AllShardsNodesResponse
	resp, err := http.Get(c.baseURL + "/shards/nodes")

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

func (c *Client) GetBalancers() (*balancerDto.AllShardsBalancersResponse, error) {
	var res balancerDto.AllShardsBalancersResponse
	resp, err := http.Get(c.baseURL + "/shards/balancers")

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
