package cache

import (
	"sync"
	"time"
)

// An item in the cache
type Item[T any] struct {
	Value      T
	Expiration int64 // Unix timestamp for expiration
	Created    int64 // Unix timestamp for creation time
}

func (item Item[T]) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Cache is an in-memory key:value store
type Cache[T any] struct {
	items map[string]Item[T]
	mu    sync.RWMutex
}

// Generic type alias for a map of items
type Items[T any] = map[string]Item[T]

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		items: make(map[string]Item[T]),
	}
}

// Set adds an item to the cache with the default expiration time
func (self *Cache[T]) Set(key string, value T) {
	self.mu.Lock()
	self.items[key] = Item[T]{
		Value:      value,
		Expiration: 0,
		Created:    time.Now().UnixNano(),
	}
	self.mu.Unlock()
}

// SetWithExpiration adds an item to the cache with a custom expiration time
func (self *Cache[T]) SetWithExpiration(key string, value T, duration time.Duration) {
	var expiration int64

	if duration <= 0 {
		// No expiry
		expiration = 0
	} else if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	self.mu.Lock()
	self.items[key] = Item[T]{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now().UnixNano(),
	}
	self.mu.Unlock()
}

// Get retrieves an item from the cache
// Returns the item and a bool indicating if the item was found
func (self *Cache[T]) Get(key string) (T, bool) {
	self.mu.RLock()
	item, found := self.items[key]
	self.mu.RUnlock()

	if !found {
		var zero T
		return zero, false
	}

	// If the item has expired, delete it and return not found
	if item.Expired() {
		self.Delete(key)
		var zero T
		return zero, false
	}

	return item.Value, true
}

// GetItem retrieves the entire cache item including metadata
func (self *Cache[T]) GetItem(key string) (Item[T], bool) {
	self.mu.RLock()
	item, found := self.items[key]
	self.mu.RUnlock()

	if !found {
		return Item[T]{}, false
	}

	if item.Expired() {
		self.Delete(key)
		return Item[T]{}, false
	}

	return item, true
}

// Delete removes an item from the cache
func (self *Cache[T]) Delete(key string) {
	self.mu.Lock()
	delete(self.items, key)
	self.mu.Unlock()
}

// Clear empties the entire cache
func (self *Cache[T]) Clear() {
	self.mu.Lock()
	self.items = make(map[string]Item[T])
	self.mu.Unlock()
}

// Count returns the number of items in the cache (including expired items)
func (self *Cache[T]) Count() int {
	self.mu.Lock()
	defer self.mu.Unlock()

	count := 0
	now := time.Now().UnixNano()

	toDelete := make([]string, 0)
	for key, item := range self.items {
		if item.Expiration == 0 || now < item.Expiration {
			count++
			continue
		}
		// Mark for deletion
		toDelete = append(toDelete, key)
	}
	// Delete items
	for _, key := range toDelete {
		delete(self.items, key)
	}

	return count
}

// ItemsNotExpired returns all non-expired items in the cache
func (self *Cache[T]) Items() Items[T] {
	self.mu.Lock()
	defer self.mu.Unlock()

	items := make(Items[T])
	now := time.Now().UnixNano()

	toDelete := make([]string, 0)
	for k, v := range self.items {
		if v.Expiration == 0 || now < v.Expiration {
			items[k] = v
			continue
		}
		// Mark for deletion
		toDelete = append(toDelete, k)
	}
	// Delete items
	for _, key := range toDelete {
		delete(self.items, key)
	}

	return items
}
