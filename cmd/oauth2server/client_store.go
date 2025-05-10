package main

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/redis/go-redis/v9"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
)

// A custom client store that persists clients in a Redis cache.
type dbClientStore struct {
	ctx   context.Context
	cache *cache.RedisCache[CacheClientInto]
}

func NewDBClientStore(ctx context.Context, cache *cache.RedisCache[CacheClientInto]) *dbClientStore {
	return &dbClientStore{
		ctx:   ctx,
		cache: cache,
	}
}

func (s *dbClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	cacheItem, err := s.cache.Get(ctx, id)
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrInvalidClient
		}
		return nil, err
	}
	return cacheItem, nil
}

func (s *dbClientStore) Set(id string, client oauth2.ClientInfo) error {
	cacheItem := CacheClientInto{
		ID:     client.GetID(),
		Secret: client.GetSecret(),
		Domain: client.GetDomain(),
		UserID: client.GetUserID(),
		Public: client.IsPublic(),
	}
	return s.cache.SetWithExpiration(s.ctx, id, cacheItem, REFRESH_TOKEN_EXP)
}

// Cachable implementation of oauth2.ClientInfo.
type CacheClientInto struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	UserID string `json:"user_id"`
	Public bool   `json:"public"`
}

func (c CacheClientInto) GetID() string {
	return c.ID
}

func (c CacheClientInto) GetSecret() string {
	return c.Secret
}

func (c CacheClientInto) GetDomain() string {
	return c.Domain
}

func (c CacheClientInto) IsPublic() bool {
	return c.Public
}

func (c CacheClientInto) GetUserID() string {
	return c.UserID
}
