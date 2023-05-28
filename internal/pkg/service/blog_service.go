package service

import (
	"github.com/so-brian/sobrian-blog-service/api"
	"github.com/so-brian/sobrian-blog-service/internal/pkg/azure/storage"
)

type BlogService interface {
	// Insert
	InsertBlog(blog api.Blog) error
	// Read
	GetBlog(name string) (*api.Blog, error)
	GetBlogs() ([]*api.Blog, error)
	// Update
	UpdateBlog(blog api.Blog) error
	// Delete
	DeleteBlog(name string) error
}

type blogService struct {
	client *storage.BlogStorageClient
}

func NewBlogService() BlogService {
	client := storage.NewAzureStorageClient()
	return &blogService{client: &client}
}

func (s *blogService) InsertBlog(blog api.Blog) error {
	return (*s.client).InsertBlog(blog)
}

func (s *blogService) GetBlog(name string) (*api.Blog, error) {
	return (*s.client).GetBlog(name)
}

func (s *blogService) GetBlogs() ([]*api.Blog, error) {
	return (*s.client).GetBlogs()
}

func (s *blogService) UpdateBlog(blog api.Blog) error {
	return (*s.client).UpdateBlog(blog)
}

func (s *blogService) DeleteBlog(name string) error {
	return (*s.client).DeleteBlog(name)
}
