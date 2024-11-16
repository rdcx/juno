package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	ctx := context.Background()

	// Start multiple containers
	for i := 0; i < 100; i++ {
		port := 7000 + i
		containerName := fmt.Sprintf("balancer_%d", i)

		// Define host paths for persistence
		hostDir := fmt.Sprintf("/home/ross/junodata/%s", containerName)
		// Ensure directories exist on the host
		os.MkdirAll(hostDir, 0755) // Creates /home/ross/junodata/balancer_i
		queueDBPath := fmt.Sprintf("%s/queue.db", hostDir)
		// if the file does not exist, create it
		if _, err := os.Stat(queueDBPath); os.IsNotExist(err) {
			file, err := os.Create(queueDBPath)
			if err != nil {
				log.Fatalf("Error creating queue.db file: %v", err)
			}
			file.Close()
		}
		policyDBPath := fmt.Sprintf("%s/policy.db", hostDir)
		// if the file does not exist, create it
		if _, err := os.Stat(policyDBPath); os.IsNotExist(err) {
			file, err := os.Create(policyDBPath)
			if err != nil {
				log.Fatalf("Error creating policy.db file: %v", err)
			}
			file.Close()
		}

		containerConfig := &container.Config{
			Image: "busybox", // Adjust the image as necessary
			Cmd:   []string{"/balancer", "-port", fmt.Sprintf("%d", port)},
			ExposedPorts: map[nat.Port]struct{}{
				nat.Port(fmt.Sprintf("%d", port) + "/tcp"): {}, // Exposing the port inside the container
			},
		}

		hostConfig := &container.HostConfig{
			NetworkMode: "host",
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%d", port) + "/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", port),
					},
				},
			},
			Binds: []string{
				"/usr/local/bin/balancer:/balancer",
				"/etc/ssl/certs:/etc/ssl/certs",
				fmt.Sprintf("%s:/queue.db", queueDBPath),
				fmt.Sprintf("%s:/policy.db", policyDBPath),
			},
		}

		resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
		if err != nil {
			log.Fatalf("Error creating container %s: %v", containerName, err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			log.Fatalf("Error starting container %s: %v", containerName, err)
		}

		log.Printf("Started container %s on port %d", containerName, port)
	}
}
