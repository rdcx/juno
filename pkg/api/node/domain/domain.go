package domain

import "github.com/google/uuid"

type Node struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"owner_id"`
	Address string    `json:"address"`
	Shards  []int     `json:"shards"`
}
