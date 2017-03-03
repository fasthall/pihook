package main

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func Pull(repo string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, repo, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}

func Run(repo string) string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, repo, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: repo,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func Copy(cid, src, dst string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, _, err := cli.CopyFromContainer(ctx, cid, src)
	if err != nil {
		panic(err)
	}
	outFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, out)
	if err != nil {
		panic(err)
	}
}
