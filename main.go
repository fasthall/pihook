package main

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

var pi string
var repo string

type PushEvent struct {
	Pusher     *Pusher     `json:"pusher" binding:"required"`
	Repository *Repository `json:"repository" binding:"required"`
	HeadCommit *HeadCommit `json:"head_commit" binding:"required"`
}

type Pusher struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type Repository struct {
	Name  string `json:"name" binding:"required"`
	Owner *Owner `json:"owner" binding:"required"`
}

type Owner struct {
	Name string `json:"name" binding:"required"`
}

type HeadCommit struct {
	ID string `json:"id" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/repo", getRepo)
	router.POST("/repo", postRepo)
	router.GET("/pi", getPi)
	router.POST("/pi", postPi)
	router.POST("/webhook", webhook)
	router.Run()
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
	fmt.Println(c.Keys)
	fmt.Println(c.Params)
	c.String(http.StatusOK, "OK")
}
