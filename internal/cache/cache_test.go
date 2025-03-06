package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// CacheSuite is a test suite for the Cache implementation
type CacheSuite struct {
	suite.Suite
	stringCache *Cache[string]
	intCache    *Cache[int]
	userCache   *Cache[User]
}

// User is a test struct for cache storage
type User struct {
	ID   int
	Name string
}

// SetupTest runs before each test
func (suite *CacheSuite) SetupTest() {
	suite.stringCache = NewCache[string]()
	suite.intCache = NewCache[int]()
	suite.userCache = NewCache[User]()
}

// TestSet tests the Set method
func (suite *CacheSuite) TestSet() {
	// String cache
	suite.stringCache.Set("key1", "value1")
	value, found := suite.stringCache.Get("key1")
	suite.True(found, "Item should be found")
	suite.Equal("value1", value, "Value should match")

	// Int cache
	suite.intCache.Set("num", 42)
	num, found := suite.intCache.Get("num")
	suite.True(found, "Item should be found")
	suite.Equal(42, num, "Value should match")

	// Struct cache
	user := User{ID: 1, Name: "John"}
	suite.userCache.Set("user1", user)
	retrievedUser, found := suite.userCache.Get("user1")
	suite.True(found, "Item should be found")
	suite.Equal(user, retrievedUser, "Value should match")
}

// TestSetWithExpiration tests the SetWithExpiration method
func (suite *CacheSuite) TestSetWithExpiration() {
	// Test with positive duration (should expire)
	suite.stringCache.SetWithExpiration("temp", "I'll expire", 10*time.Millisecond)

	// Verify it exists initially
	value, found := suite.stringCache.Get("temp")
	suite.True(found, "Item should be found before expiration")
	suite.Equal("I'll expire", value, "Value should match")

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Verify it's gone after expiration
	_, found = suite.stringCache.Get("temp")
	suite.False(found, "Item should not be found after expiration")

	// Test with zero duration (should never expire)
	suite.stringCache.SetWithExpiration("forever", "I won't expire", 0)

	// Wait a bit
	time.Sleep(20 * time.Millisecond)

	// Verify it still exists
	value, found = suite.stringCache.Get("forever")
	suite.True(found, "Item with zero expiration should not expire")
	suite.Equal("I won't expire", value, "Value should match")

	// Test with negative duration (should never expire)
	suite.stringCache.SetWithExpiration("also-forever", "I won't expire either", -1*time.Second)

	// Wait a bit
	time.Sleep(20 * time.Millisecond)

	// Verify it still exists
	value, found = suite.stringCache.Get("also-forever")
	suite.True(found, "Item with negative expiration should not expire")
	suite.Equal("I won't expire either", value, "Value should match")
}

// TestGetItem tests the GetItem method
func (suite *CacheSuite) TestGetItem() {
	now := time.Now().UnixNano()

	suite.stringCache.Set("key1", "value1")

	item, found := suite.stringCache.GetItem("key1")
	suite.True(found, "Item should be found")
	suite.Equal("value1", item.Value, "Value should match")
	suite.Equal(int64(0), item.Expiration, "No expiration should be 0")
	suite.True(item.Created >= now, "Created timestamp should be after test start")

	// Test with expiration
	suite.stringCache.SetWithExpiration("temp", "expiring", 100*time.Millisecond)

	item, found = suite.stringCache.GetItem("temp")
	suite.True(found, "Item should be found")
	suite.Equal("expiring", item.Value, "Value should match")
	suite.True(item.Expiration > now, "Expiration should be in the future")

	// Test expired item
	suite.stringCache.SetWithExpiration("expired", "gone", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, found = suite.stringCache.GetItem("expired")
	suite.False(found, "Expired item should not be found")
}

// TestDelete tests the Delete method
func (suite *CacheSuite) TestDelete() {
	// Add an item
	suite.stringCache.Set("key1", "value1")

	// Verify it exists
	_, found := suite.stringCache.Get("key1")
	suite.True(found, "Item should be found before deletion")

	// Delete it
	suite.stringCache.Delete("key1")

	// Verify it's gone
	_, found = suite.stringCache.Get("key1")
	suite.False(found, "Item should not be found after deletion")

	// Delete non-existent key (should not panic)
	suite.NotPanics(func() {
		suite.stringCache.Delete("non-existent")
	}, "Deleting non-existent key should not panic")
}

// TestClear tests the Clear method
func (suite *CacheSuite) TestClear() {
	// Add some items
	suite.stringCache.Set("key1", "value1")
	suite.stringCache.Set("key2", "value2")

	// Verify they exist
	count := suite.stringCache.Count()
	suite.Equal(2, count, "Cache should have 2 items before clearing")

	// Clear the cache
	suite.stringCache.Clear()

	// Verify it's empty
	count = suite.stringCache.Count()
	suite.Equal(0, count, "Cache should be empty after clearing")
}

// TestCount tests the Count method
func (suite *CacheSuite) TestCount() {
	// Empty cache
	count := suite.stringCache.Count()
	suite.Equal(0, count, "Empty cache should have 0 items")

	// Add some items
	suite.stringCache.Set("key1", "value1")
	suite.stringCache.Set("key2", "value2")

	count = suite.stringCache.Count()
	suite.Equal(2, count, "Cache should have 2 items")

	// Add an item that expires quickly
	suite.stringCache.SetWithExpiration("temp", "expiring", 10*time.Millisecond)

	count = suite.stringCache.Count()
	suite.Equal(3, count, "Cache should have 3 items before expiration")

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Count should automatically clean up expired items
	count = suite.stringCache.Count()
	suite.Equal(2, count, "Count should remove expired items")
}

// TestItems tests the Items method
func (suite *CacheSuite) TestItems() {
	// Add some items
	suite.stringCache.Set("key1", "value1")
	suite.stringCache.Set("key2", "value2")
	suite.stringCache.SetWithExpiration("temp", "expiring", 10*time.Millisecond)

	// Get all items
	items := suite.stringCache.Items()
	suite.Equal(3, len(items), "Items should return all 3 items before expiration")

	// Verify item values
	suite.Equal("value1", items["key1"].Value, "Item value should match")
	suite.Equal("value2", items["key2"].Value, "Item value should match")
	suite.Equal("expiring", items["temp"].Value, "Item value should match")

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Get items again - should clean up expired
	items = suite.stringCache.Items()
	suite.Equal(2, len(items), "Items should clean up expired items")
	suite.Contains(items, "key1", "Items should contain permanent keys")
	suite.Contains(items, "key2", "Items should contain permanent keys")
	suite.NotContains(items, "temp", "Items should not contain expired keys")
}

// TestConcurrentAccess tests concurrent access to the cache
func (suite *CacheSuite) TestConcurrentAccess() {
	const goroutines = 10
	const operationsPerGoroutine = 100

	done := make(chan bool, goroutines*2)

	// Start writers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				suite.stringCache.Set(key, fmt.Sprintf("value-%d-%d", id, j))
			}
			done <- true
		}(i)
	}

	// Start readers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				// Mix of operations: get, count, items
				switch j % 3 {
				case 0:
					key := fmt.Sprintf("key-%d-%d", id, j)
					suite.stringCache.Get(key)
				case 1:
					suite.stringCache.Count()
				case 2:
					suite.stringCache.Items()
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < goroutines*2; i++ {
		<-done
	}

	// If we got here without deadlock or panic, the test passes
}

// TestZeroValueHandling tests that zero values are handled correctly
func (suite *CacheSuite) TestZeroValueHandling() {
	// Initialize with zero values
	suite.intCache.Set("zero", 0)
	suite.stringCache.Set("empty", "")
	suite.userCache.Set("nobody", User{})

	// Test getting the zero values
	zeroInt, found := suite.intCache.Get("zero")
	suite.True(found, "Zero int should be found")
	suite.Equal(0, zeroInt, "Zero int value should be 0")

	emptyString, found := suite.stringCache.Get("empty")
	suite.True(found, "Empty string should be found")
	suite.Equal("", emptyString, "Empty string value should be empty")

	emptyUser, found := suite.userCache.Get("nobody")
	suite.True(found, "Empty user should be found")
	suite.Equal(User{}, emptyUser, "Empty user should equal the zero value for User")

	// Test missing values
	missingInt, found := suite.intCache.Get("missing")
	suite.False(found, "Missing key should not be found")
	suite.Equal(0, missingInt, "Missing int should return zero value (0)")

	missingString, found := suite.stringCache.Get("missing")
	suite.False(found, "Missing key should not be found")
	suite.Equal("", missingString, "Missing string should return zero value (empty string)")

	missingUser, found := suite.userCache.Get("missing")
	suite.False(found, "Missing key should not be found")
	suite.Equal(User{}, missingUser, "Missing user should return zero value for User")
}

// TestItemExpired tests the Expired method of Item directly
func (suite *CacheSuite) TestItemExpired() {
	now := time.Now()

	tests := []struct {
		name       string
		expiration int64
		expected   bool
	}{
		{
			name:       "zero expiration",
			expiration: 0,
			expected:   false,
		},
		{
			name:       "future expiration",
			expiration: now.Add(1 * time.Minute).UnixNano(),
			expected:   false,
		},
		{
			name:       "past expiration",
			expiration: now.Add(-1 * time.Minute).UnixNano(),
			expected:   true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			item := Item[string]{
				Value:      "test",
				Expiration: tc.expiration,
				Created:    now.UnixNano(),
			}

			suite.Equal(tc.expected, item.Expired(), "Expired() should return expected value")
		})
	}
}

// TestNewCache tests the NewCache constructor
func (suite *CacheSuite) TestNewCache() {
	// Test with different types
	stringCache := NewCache[string]()
	suite.NotNil(stringCache, "NewCache should return a non-nil cache")
	suite.NotNil(stringCache.items, "Cache items map should be initialized")

	intCache := NewCache[int]()
	suite.NotNil(intCache, "NewCache should return a non-nil cache")

	// Custom type
	type Person struct {
		Name string
		Age  int
	}

	personCache := NewCache[Person]()
	suite.NotNil(personCache, "NewCache should return a non-nil cache")
}

// Run the test suite
func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
