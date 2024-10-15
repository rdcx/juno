package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"math/rand"
)

func randomShardRange() string {
	rand.Seed(time.Now().UnixNano())

	// Define the maximum offset and size limits
	const maxOffset = 100000
	const maxSize = 500

	// Generate a random offset within the range [0, 100000 - maxSize]
	offset := rand.Intn(maxOffset - maxSize + 1)

	// Generate a random size within the range [1, maxSize]
	size := rand.Intn(maxSize) + 1

	// Return the shard range as a string in the format "[offset,length]"
	return fmt.Sprintf("[%d,%d]", offset, size)
}

func main() {
	nodes := 10000

	for i := 0; i < nodes; i++ {
		req, err := http.NewRequest("POST", "http://localhost:8080/nodes", bytes.NewBuffer(
			[]byte(fmt.Sprintf(`{"address":"node%d.com:9392","shard_assignments":[%s]}`, i, randomShardRange())),
		))

		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJvc3NAZXhhbXBsZS5jb20iLCJleHAiOjE3MjkwMjQ0MDMsImlkIjoiZjA2MDJiYTMtMzBlMy00ODQ5LTkyMzQtNTMxNWQ0OGU3OTFjIn0.UvobuhQooezh0x1zGtw_vBAk7YLOnhNZQZiqFZqaYk4")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusCreated {
			panic("unexpected status code")
		}

		resp.Body.Close()
	}

	fmt.Printf("Created %d nodes\n", nodes)
}
