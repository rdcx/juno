package main

import (
	"context"
	"fmt"
	"log"

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
		port := 6000 + i
		containerName := fmt.Sprintf("ranag_%d", i)

		containerConfig := &container.Config{
			Image: "busybox", // Adjust the image as necessary
			Cmd:   []string{"/ranag", "-port", fmt.Sprintf("%d", port)},
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
				"/usr/local/bin/ranag:/ranag",
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
