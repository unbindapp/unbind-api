package main

import (
	"context"
	"sync"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
)

type clientStore struct {
	mu      sync.RWMutex
	clients map[string]CacheClientInto
}

func NewClientStore() *clientStore {
	return &clientStore{
		clients: make(map[string]CacheClientInto),
	}
}

func (s *clientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, ok := s.clients[id]
	if !ok {
		return nil, errors.ErrInvalidClient
	}
	return client, nil
}

func (s *clientStore) Set(id string, client oauth2.ClientInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[id] = CacheClientInto{
		ID:     client.GetID(),
		Secret: client.GetSecret(),
		Domain: client.GetDomain(),
		UserID: client.GetUserID(),
		Public: client.IsPublic(),
	}
}

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
