package repository

import (
	"clase03-memcached/internal/domain"
	"context"
	"fmt"
	"github.com/karlseguin/ccache"
	"time"
)

type ItemsLocalCacheRepository struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewItemsLocalCacheRepository(ttl time.Duration) *ItemsLocalCacheRepository {
	return &ItemsLocalCacheRepository{
		client: ccache.New(ccache.Configure[string]()),
		ttl:    ttl,
	}
}

func (r ItemsLocalCacheRepository) List(ctx context.Context) ([]domain.Item, error) {
	return nil, fmt.Errorf("list is not supported in memcached")
}

func (r ItemsLocalCacheRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	r.client.Set(item.ID, item, r.ttl)
	return item, nil
}

func (r ItemsLocalCacheRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	it := r.client.Get(id)
	if it == nil {
		return domain.Item{}, fmt.Errorf("item not found in cache")
	}
	item, ok := it.Value().(domain.Item)
	if !ok {
		return domain.Item{}, fmt.Errorf("error asserting item type from cache")
	}
	return item, nil
}

func (r ItemsLocalCacheRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (r ItemsLocalCacheRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
