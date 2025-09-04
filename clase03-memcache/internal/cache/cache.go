package cache

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
	"example.com/gin-memcached-base/internal/models"
)

type Client struct {
	mc *memcache.Client
}

func NewClient(addr string) *Client {
	return &Client{
		mc: memcache.New(addr), // ejemplo "127.0.0.1:11211"
	}
}

// Guardar un item individual
func (c *Client) Set(key string, item models.Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return c.mc.Set(&memcache.Item{Key: key, Value: data})
}

func (c *Client) Get(key string) (models.Item, error) {
	it, err := c.mc.Get(key)
	if err != nil {
		return models.Item{}, err
	}
	var item models.Item
	if err := json.Unmarshal(it.Value, &item); err != nil {
		return models.Item{}, err
	}
	return item, nil
}

// Guardar lista de items
func (c *Client) SetList(key string, items []models.Item) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return c.mc.Set(&memcache.Item{Key: key, Value: data})
}

func (c *Client) GetList(key string) ([]models.Item, error) {
	it, err := c.mc.Get(key)
	if err != nil {
		return nil, err
	}
	var items []models.Item
	if err := json.Unmarshal(it.Value, &items); err != nil {
		return nil, err
	}
	return items, nil
}
