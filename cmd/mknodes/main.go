package main

import (
	"bytes"
	"fmt"
	"io"
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

func port(i int) string {
	if i > 99 {
		return fmt.Sprintf("9%d", i)
	}

	if i > 9 {
		return fmt.Sprintf("90%d", i)
	}

	return fmt.Sprintf("900%d", i)
}

func shardRange(i int) string {
	// make the shard range 100k / 100
	offset := i * 1000
	size := 1000

	return fmt.Sprintf("[%d,%d]\n", offset, size)
}

func main() {
	nodes := 100

	for i := 0; i < nodes; i++ {
		req, err := http.NewRequest("POST", "http://localhost:8080/nodes", bytes.NewBuffer(
			[]byte(fmt.Sprintf(`{"address":"127.0.0.1:%s","shard_assignments":[%s]}`, port(i), shardRange(i))),
		))

		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJvc3NAZXhhbXBsZS5jb20iLCJleHAiOjE3MzAxMTg1MDQsImlkIjoiNDMwMzE2NzYtYjgxZC00ODE5LTg3MjktNTJiNWY1MzQ0MTViIiwibmFtZSI6IlJvc3MifQ.NIDuuGwW3eejo8yJ28rdgrw93SsguO0OaHv_N85e8iI")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusCreated {
			fmt.Println(io.ReadAll(resp.Body))
			panic("unexpected status code")
		}

		resp.Body.Close()
	}

	fmt.Printf("Created %d nodes\n", nodes)
}
