package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
)

var pi string
var repo string

func main() {
	router := gin.Default()
	router.GET("/repo", getRepo)
	router.POST("/repo", postRepo)
	router.GET("/pi", getPi)
	router.POST("/pi", postPi)
	router.POST("/webhook", webhook)
	router.GET("/test", test)
	router.Run()
}

func test(c *gin.Context) {
	u, err := uuid.NewV4()
	cmd := exec.Command("git", "clone", repo, u.String())
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Couldn't clone the repo %s", repo))
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cloned into", path.Join(dir, u.String()))
	cid := runContainer("fasthall/smartfarm_sketch", path.Join(dir, u.String()))
	fmt.Println(cid)
}

func getRepo(c *gin.Context) {
	c.String(http.StatusOK, repo)
}

func postRepo(c *gin.Context) {
	repo = c.Query("host")
	c.String(http.StatusOK, c.Query("host")+" added\n")
}

func getPi(c *gin.Context) {
	c.String(http.StatusOK, pi)
}

func postPi(c *gin.Context) {
	pi = c.Query("host")
	c.String(http.StatusOK, c.Query("host")+" added\n")
}

func webhook(c *gin.Context) {
	b := []byte{}
	c.Request.Body.Read(b)
	event := c.Request.Header.Get("X-GitHub-Event")
	if event == "push" {
		os.Chdir(path.Join(os.Getenv("HOME"), "smartfarm_sketch"))
		cmd := "git"
		args := []string{"pull"}
		out, err := exec.Command(cmd, args...).Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	} else if event == "ping" {
		fmt.Println("Github is testing!")
	}
	c.String(http.StatusOK, "OK")
}

func runContainer(repo, bind string) string {
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
	}, &container.HostConfig{
		AutoRemove: true,
		Binds:      []string{bind + ":/smartfarm_sketch/sketch/"},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}
