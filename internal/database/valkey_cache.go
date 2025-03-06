package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

// ValueCoder handles encoding and decoding of cache values.
// This allows for generic type support with Valkey.
type ValueCoder[T any] interface {
	Encode(T) (string, error)
	Decode(string) (T, error)
}

// JSONValueCoder implements ValueCoder using JSON encoding/decoding.
type JSONValueCoder[T any] struct{}

// Encode converts a value to a JSON string.
func (c JSONValueCoder[T]) Encode(value T) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to JSON encode: %w", err)
	}
	return string(data), nil
}

// Decode converts a JSON string back to a value.
func (c JSONValueCoder[T]) Decode(data string) (T, error) {
	var value T
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return value, fmt.Errorf("failed to JSON decode: %w", err)
	}
	return value, nil
}

// StringValueCoder is a specialized coder for string values that doesn't need JSON.
type StringValueCoder struct{}

// Encode simply returns the string as is.
func (c StringValueCoder) Encode(value string) (string, error) {
	return value, nil
}

// Decode simply returns the string as is.
func (c StringValueCoder) Decode(data string) (string, error) {
	return data, nil
}

// ValkeyCache is a wrapper around a Valkey client that adds generic type support.
type ValkeyCache[T any] struct {
	client    valkey.Client
	keyPrefix string
	coder     ValueCoder[T]
}

// NewCache creates a new cache with the specified Valkey client.
func NewCache[T any](client valkey.Client, keyPrefix string) *ValkeyCache[T] {
	return &ValkeyCache[T]{
		client:    client,
		keyPrefix: keyPrefix,
		coder:     JSONValueCoder[T]{},
	}
}

// NewStringCache creates a specialized cache for string values that avoids JSON overhead.
func NewStringCache(client valkey.Client, keyPrefix string) *ValkeyCache[string] {
	return &ValkeyCache[string]{
		client:    client,
		keyPrefix: keyPrefix,
		coder:     StringValueCoder{},
	}
}

// WithCoder allows customizing the value encoding/decoding mechanism.
func (self *ValkeyCache[T]) WithCoder(coder ValueCoder[T]) *ValkeyCache[T] {
	self.coder = coder
	return self
}

// fullKey generates the full key with prefix.
func (self *ValkeyCache[T]) fullKey(key string) string {
	if self.keyPrefix == "" {
		return key
	}
	return self.keyPrefix + ":" + key
}

// Set adds an item to the cache with no expiration.
func (self *ValkeyCache[T]) Set(ctx context.Context, key string, value T) error {
	encoded, err := self.coder.Encode(value)
	if err != nil {
		return err
	}

	cmd := self.client.B().Set().Key(self.fullKey(key)).Value(encoded).Build()
	return self.client.Do(ctx, cmd).Error()
}

// SetWithExpiration adds an item to the cache with an expiration time.
func (self *ValkeyCache[T]) SetWithExpiration(ctx context.Context, key string, value T, expiration time.Duration) error {
	encoded, err := self.coder.Encode(value)
	if err != nil {
		return err
	}

	cmd := self.client.B().Set().Key(self.fullKey(key)).Value(encoded).Ex(expiration).Build()
	return self.client.Do(ctx, cmd).Error()
}

// ErrKeyNotFound is returned when a key is not found in the cache.
var ErrKeyNotFound = errors.New("key not found in cache")

// Get retrieves an item from the cache.
func (c *ValkeyCache[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	cmd := c.client.B().Get().Key(c.fullKey(key)).Build()
	result, err := c.client.Do(ctx, cmd).ToString()
	if err != nil {
		if err == valkey.Nil {
			return value, ErrKeyNotFound
		}
		return value, err
	}

	return c.coder.Decode(result)
}

// GetWithTTL retrieves an item and its remaining TTL from the cache.
func (c *ValkeyCache[T]) GetWithTTL(ctx context.Context, key string) (T, time.Duration, error) {
	var value T
	fullKey := c.fullKey(key)

	// First check if the key exists
	getCmd := c.client.B().Get().Key(fullKey).Build()
	result, err := c.client.Do(ctx, getCmd).ToString()
	if err != nil {
		if err == valkey.Nil {
			return value, 0, ErrKeyNotFound
		}
		return value, 0, err
	}

	// Get TTL
	ttlCmd := c.client.B().Ttl().Key(fullKey).Build()
	ttlSeconds, err := c.client.Do(ctx, ttlCmd).ToInt64()
	if err != nil {
		return value, 0, err
	}

	// Convert TTL to duration
	var ttl time.Duration
	if ttlSeconds > 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	} else if ttlSeconds == -1 {
		// Key exists but has no expiration
		ttl = -1
	} else {
		// Key doesn't exist or is about to expire
		return value, 0, ErrKeyNotFound
	}

	// Decode value
	value, err = c.coder.Decode(result)
	if err != nil {
		return value, 0, err
	}

	return value, ttl, nil
}

// Delete removes an item from the cache.
func (c *ValkeyCache[T]) Delete(ctx context.Context, key string) error {
	cmd := c.client.B().Del().Key(c.fullKey(key)).Build()
	return c.client.Do(ctx, cmd).Error()
}

// DeleteAll removes all items with the cache's prefix.
func (c *ValkeyCache[T]) DeleteAll(ctx context.Context) error {
	if c.keyPrefix == "" {
		return errors.New("cannot delete all keys without a prefix")
	}

	// Find all keys with this prefix
	pattern := c.keyPrefix + ":*"
	keysCmd := c.client.B().Keys().Pattern(pattern).Build()
	keys, err := c.client.Do(ctx, keysCmd).AsStrSlice()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all found keys
	delCmd := c.client.B().Del().Key(keys...).Build()
	return c.client.Do(ctx, delCmd).Error()
}

// Exists checks if a key exists in the cache.
func (c *ValkeyCache[T]) Exists(ctx context.Context, key string) (bool, error) {
	cmd := c.client.B().Exists().Key(c.fullKey(key)).Build()
	result, err := c.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Keys returns all keys with the cache's prefix.
func (c *ValkeyCache[T]) Keys(ctx context.Context) ([]string, error) {
	var pattern string
	if c.keyPrefix == "" {
		pattern = "*"
	} else {
		pattern = c.keyPrefix + ":*"
	}

	cmd := c.client.B().Keys().Pattern(pattern).Build()
	keys, err := c.client.Do(ctx, cmd).AsStrSlice()
	if err != nil {
		return nil, err
	}

	// Remove prefix from keys for consistency
	if c.keyPrefix != "" {
		prefixLen := len(c.keyPrefix) + 1 // +1 for the colon
		for i, key := range keys {
			keys[i] = key[prefixLen:]
		}
	}

	return keys, nil
}

// GetAll retrieves all items from the cache.
func (c *ValkeyCache[T]) GetAll(ctx context.Context) (map[string]T, error) {
	keys, err := c.Keys(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]T)
	for _, key := range keys {
		value, err := c.Get(ctx, key)
		if err != nil && !errors.Is(err, ErrKeyNotFound) {
			return nil, err
		}
		if !errors.Is(err, ErrKeyNotFound) {
			result[key] = value
		}
	}

	return result, nil
}

// Increment increments a numeric value stored at the given key.
// Only works for numeric types like int, int64, float64.
func Increment[T ~int | ~int64 | ~float64](
	ctx context.Context,
	cache *ValkeyCache[T],
	key string,
	increment T) (T, error) {

	var result T
	fullKey := cache.fullKey(key)

	// Use INCRBY for integers and INCRBYFLOAT for floats
	var cmd valkey.Completed
	switch any(increment).(type) {
	case int, int64:
		cmd = cache.client.B().Incrby().Key(fullKey).Increment(int64(increment)).Build()
		val, err := cache.client.Do(ctx, cmd).ToInt64()
		if err != nil {
			return result, err
		}
		return T(val), nil
	case float64:
		cmd = cache.client.B().Incrbyfloat().Key(fullKey).Increment(float64(increment)).Build()
		val, err := cache.client.Do(ctx, cmd).ToFloat64()
		if err != nil {
			return result, err
		}
		return T(val), nil
	default:
		return result, fmt.Errorf("unsupported type for increment")
	}
}
