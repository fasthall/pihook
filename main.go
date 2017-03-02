package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/gin-gonic/gin"
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
	os.Chdir(path.Join(os.Getenv("HOME"), "smartfarm_sketch"))
	cmd := "git"
	args := []string{"pull"}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
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
