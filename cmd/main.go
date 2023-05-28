package main

import (
	"github.com/gin-gonic/gin"
	"github.com/so-brian/sobrian-blog-service/api"
	"github.com/so-brian/sobrian-blog-service/internal/pkg/azure/storage"
	"github.com/so-brian/sobrian-blog-service/internal/pkg/server"
)

func main() {
	println("Hello World")

	client := storage.NewAzureStorageClient()
	client.InsertBlog(api.Blog{
		Title: "Hello World2",
		Body:  "This is a test2",
	})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server := server.NewBlogServer()
	r.POST("/blog", server.InsertBlog)
	r.GET("/blog/:name", server.GetBlog)
	r.GET("/blogs", server.GetBlogs)
	r.PUT("/blog", server.UpdateBlog)
	r.DELETE("/blog/:name", server.DeleteBlog)

	r.Run()
}
