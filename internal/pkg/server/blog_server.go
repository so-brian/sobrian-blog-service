package server

import (
	"github.com/gin-gonic/gin"
	"github.com/so-brian/sobrian-blog-service/api"
	"github.com/so-brian/sobrian-blog-service/internal/pkg/service"
)

type BlogServer interface {
	// Insert
	InsertBlog(*gin.Context)
	// Read
	GetBlog(*gin.Context)
	GetBlogs(*gin.Context)
	// Update
	UpdateBlog(*gin.Context)
	// Delete
	DeleteBlog(*gin.Context)
}

type blogServer struct {
	service *service.BlogService
}

func NewBlogServer() BlogServer {
	service := service.NewBlogService()
	return &blogServer{service: &service}
}

func (s *blogServer) InsertBlog(c *gin.Context) {
	var blog api.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := (*s.service).InsertBlog(blog); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func (s *blogServer) GetBlog(c *gin.Context) {
	name := c.Param("name")
	blog, err := (*s.service).GetBlog(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blog)
}

func (s *blogServer) GetBlogs(c *gin.Context) {
	blogs, err := (*s.service).GetBlogs()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogs)
}

func (s *blogServer) UpdateBlog(c *gin.Context) {
	var blog api.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := (*s.service).UpdateBlog(blog); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func (s *blogServer) DeleteBlog(c *gin.Context) {
	name := c.Param("name")
	if err := (*s.service).DeleteBlog(name); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
