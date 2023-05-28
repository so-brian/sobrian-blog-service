package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/so-brian/sobrian-blog-service/api"
)

const (
	containerName = "blog"
)

type BlogStorageClient interface {
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

type azureStorageClient struct {
	client *azblob.Client
}

func NewAzureStorageClient() BlogStorageClient {
	// https://learn.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go?tabs=roles-azure-portal
	// TODO: replace <storage-account-name> with your actual storage account name
	url := "https://sobrian.blob.core.windows.net/"
	// ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	return &azureStorageClient{client: client}
}

func (c *azureStorageClient) InsertBlog(blog api.Blog) error {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(blog)
	res, err := c.client.UploadStream(context.Background(), containerName, fmt.Sprint(blog.Id), &buf, &azblob.UploadStreamOptions{})
	if err != nil {
		log.Println(err.Error())

		return err
	}

	log.Println(res.ContentMD5)
	return nil
}

func (c *azureStorageClient) GetBlog(name string) (*api.Blog, error) {
	// Download the blob
	ctx := context.Background()
	res, err := c.client.DownloadStream(ctx, containerName, name, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Read the blob
	downloadedData := bytes.Buffer{}
	retryReader := res.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Close the blob
	err = retryReader.Close()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Decode the blob
	var blog api.Blog
	err = json.NewDecoder(&downloadedData).Decode(&blog)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &blog, nil
}

func (c *azureStorageClient) GetBlogs() ([]*api.Blog, error) {
	// blob listings are returned across multiple pages
	pager := c.client.NewListBlobsFlatPager(containerName, nil)

	var blogs []*api.Blog

	// continue fetching pages until no more remain
	for pager.More() {
		// advance to the next page
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		// print the blob names for this page
		for _, blob := range page.Segment.BlobItems {
			blog, err := c.GetBlog(*blob.Name)
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}

			blogs = append(blogs, blog)
		}
	}

	return blogs, nil
}

func (c *azureStorageClient) UpdateBlog(blog api.Blog) error {
	return c.InsertBlog(blog)
}

func (c *azureStorageClient) DeleteBlog(name string) error {
	_, err := c.client.DeleteBlob(context.Background(), containerName, name, nil)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}
