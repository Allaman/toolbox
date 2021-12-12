package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	CONTAINER_NAME = "debian"
	IMAGE_NAME     = "debian:latest"
)

func removeContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}
	if err := cli.ContainerRemove(ctx, CONTAINER_NAME, removeOptions); err != nil {
		fmt.Printf("Unable to remove container: %s", err)
		panic(err)
	}
	fmt.Println("container removed")
}

func checkForRunningContainer() bool {
	var is_running bool
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], CONTAINER_NAME) && container.State == "running" {
			is_running = true
		}
		if strings.Contains(container.Names[0], CONTAINER_NAME) && container.State == "exited" {
			removeContainer()
		}
	}
	return is_running
}

func startContainer() string {
	if checkForRunningContainer() {
		fmt.Println("container is already running")
		return ""
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, IMAGE_NAME, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: IMAGE_NAME,
		Cmd:   []string{"echo", "hello world"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"1022/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "1022",
				},
			},
		},
	}, nil, nil, CONTAINER_NAME)
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("container started")
	return resp.ID
}

func getContainerLogs(id string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	options := types.ContainerLogsOptions{ShowStdout: true}
	out, err := cli.ContainerLogs(ctx, id, options)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}
