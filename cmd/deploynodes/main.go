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
		port := 9000 + i
		containerName := fmt.Sprintf("node_%d", i)

		// Define host paths for persistence
		hostDir := fmt.Sprintf("/home/ross/junodata/%s", containerName)
		storagePath := fmt.Sprintf("%s/storage", hostDir)
		// Ensure directories exist on the host
		os.MkdirAll(storagePath, 0755) // Creates /home/ross/junodata/node_i/storage
		pageDBPath := fmt.Sprintf("%s/page.db", hostDir)
		// if the file does not exist, create it
		if _, err := os.Stat(pageDBPath); os.IsNotExist(err) {
			file, err := os.Create(pageDBPath)
			if err != nil {
				log.Fatalf("Error creating page.db file: %v", err)
			}
			file.Close()
		}

		containerConfig := &container.Config{
			Image: "busybox", // Adjust the image as necessary
			Cmd:   []string{"/node", "-port", fmt.Sprintf("%d", port)},
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
				"/usr/local/bin/node:/node",
				"/etc/ssl/certs:/etc/ssl/certs",
				fmt.Sprintf("%s:/page.db", pageDBPath),
				fmt.Sprintf("%s:/storage", storagePath),
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
