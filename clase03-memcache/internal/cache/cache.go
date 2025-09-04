package cache

import (
	"errors"

	"example.com/gin-memcached-base/internal/models"
)

type Client struct {
	storage map[string]models.Item
}

func NewClient() *Client {
	return &Client{
		storage: make(map[string]models.Item),
	}
}

func (c *Client) Set(key string, item models.Item) error {
	c.storage[key] = item
	return nil
}

func (c *Client) Get(key string) (models.Item, error) {
	item, exists := c.storage[key]
	if !exists {
		return models.Item{}, errors.New("item no encontrado")
	}
	return item, nil
}