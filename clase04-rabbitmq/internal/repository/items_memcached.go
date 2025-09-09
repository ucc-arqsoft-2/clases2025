package repository

import (
	"clase04-rabbitmq/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type MemcachedItemsRepository struct {
	ttl    time.Duration
	client *memcache.Client
}

func NewMemcachedItemsRepository(host string, port string, ttl time.Duration) MemcachedItemsRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	return MemcachedItemsRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r MemcachedItemsRepository) List(ctx context.Context) ([]domain.Item, error) {
	return nil, fmt.Errorf("list is not supported in memcached")
}

func (r MemcachedItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error marshalling item to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        item.ID,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return domain.Item{}, fmt.Errorf("error setting item in memcached: %w", err)
	}
	return item, nil
}

func (r MemcachedItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	bytes, err := r.client.Get(id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error getting item from memcached: %w", err)
	}
	var item domain.Item
	if err := json.Unmarshal(bytes.Value, &item); err != nil {
		return domain.Item{}, fmt.Errorf("error unmarshalling item from JSON: %w", err)
	}
	return item, nil
}

func (r MemcachedItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (r MemcachedItemsRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
