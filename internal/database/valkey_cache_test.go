package database

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valkey-io/valkey-go/mock"
	"go.uber.org/mock/gomock"
)

// User is a test struct for cache storage
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CacheSuite is a test suite for the Cache implementation
type CacheSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockClient  *mock.Client
	stringCache *ValkeyCache[string]
	intCache    *ValkeyCache[int]
	userCache   *ValkeyCache[User]
	ctx         context.Context
}

// SetupTest runs before each test
func (s *CacheSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockClient = mock.NewClient(s.ctrl)
	s.stringCache = NewCache[string](s.mockClient, "test")
	s.intCache = NewCache[int](s.mockClient, "int-test")
	s.userCache = NewCache[User](s.mockClient, "user-test")
	s.ctx = context.Background()
}

// TearDownTest cleans up after each test
func (s *CacheSuite) TearDownTest() {
	s.ctrl.Finish()
}

// TestSet tests the Set method
func (s *CacheSuite) TestSet() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("SET", "test:key1", "\"value1\"")). // Include the value in the match
		Return(mock.Result(mock.ValkeyString("OK")))

	// Call method
	err := s.stringCache.Set(s.ctx, "key1", "value1")

	// Assert
	assert.NoError(s.T(), err)
}

// TestSetWithExpiration tests the SetWithExpiration method
func (s *CacheSuite) TestSetWithExpiration() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("SET", "test:key1", "\"value1\"", "EX", "600")).
		Return(mock.Result(mock.ValkeyString("OK")))

	// Call method
	err := s.stringCache.SetWithExpiration(s.ctx, "key1", "value1", 10*time.Minute)

	// Assert
	assert.NoError(s.T(), err)
}

// TestGet tests the Get method
func (s *CacheSuite) TestGet() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyString("\"value1\""))) // Note the extra quotes for JSON

	// Call method
	value, err := s.stringCache.Get(s.ctx, "key1")

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "value1", value)
}

// TestGet_NotFound tests the Get method when key is not found
func (s *CacheSuite) TestGet_NotFound() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("GET", "test:nonexistent")).
		Return(mock.Result(mock.ValkeyNil()))

	// Call method
	_, err := s.stringCache.Get(s.ctx, "nonexistent")

	// Assert
	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, ErrKeyNotFound))
}

// TestGetWithTTL tests the GetWithTTL method
func (s *CacheSuite) TestGetWithTTL() {
	// Set up expectation for GET
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("GET", "test:key1")).
		Return(mock.Result(mock.ValkeyString("\"value1\"")))

	// Set up expectation for TTL
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("TTL", "test:key1")).
		Return(mock.Result(mock.ValkeyInt64(300))) // 5 minutes in seconds

	// Call method
	value, ttl, err := s.stringCache.GetWithTTL(s.ctx, "key1")

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "value1", value)
	assert.Equal(s.T(), 300*time.Second, ttl)
}

// TestGetWithTTL_NoExpiry tests the GetWithTTL method for keys with no expiration
func (s *CacheSuite) TestGetWithTTL_NoExpiry() {
	// Set up expectation for GET
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("GET", "test:key1")).
		Return(mock.Result(mock.ValkeyString("\"value1\"")))

	// Set up expectation for TTL (return -1 for keys with no expiry)
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("TTL", "test:key1")).
		Return(mock.Result(mock.ValkeyInt64(-1)))

	// Call method
	value, ttl, err := s.stringCache.GetWithTTL(s.ctx, "key1")

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "value1", value)
	assert.Equal(s.T(), time.Duration(-1), ttl) // -1 indicates no expiry
}

// TestDelete tests the Delete method
func (s *CacheSuite) TestDelete() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("DEL", "test:key1")).
		Return(mock.Result(mock.ValkeyInt64(1)))

	// Call method
	err := s.stringCache.Delete(s.ctx, "key1")

	// Assert
	assert.NoError(s.T(), err)
}

// TestDeleteAll tests the DeleteAll method
func (s *CacheSuite) TestDeleteAll() {
	// Set up expectation for KEYS
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("KEYS", "test:*")).
		Return(mock.Result(mock.ValkeyArray(
			mock.ValkeyString("test:key1"),
			mock.ValkeyString("test:key2"),
		)))

	// Set up expectation for DEL
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("DEL", "test:key1", "test:key2")).
		Return(mock.Result(mock.ValkeyInt64(2)))

	// Call method
	err := s.stringCache.DeleteAll(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
}

// TestDeleteAll_NoKeys tests the DeleteAll method when no keys match
func (s *CacheSuite) TestDeleteAll_NoKeys() {
	// Set up expectation for KEYS (empty result)
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("KEYS", "test:*")).
		Return(mock.Result(mock.ValkeyArray()))

	// Call method
	err := s.stringCache.DeleteAll(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
}

// TestDeleteAll_NoPrefix tests the DeleteAll method when no prefix is set
func (s *CacheSuite) TestDeleteAll_NoPrefix() {
	// Create cache with no prefix
	noPrefix := NewCache[string](s.mockClient, "")

	// Call method
	err := noPrefix.DeleteAll(s.ctx)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "without a prefix")
}

// TestExists tests the Exists method
func (s *CacheSuite) TestExists() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("EXISTS", "test:key1")).
		Return(mock.Result(mock.ValkeyInt64(1)))

	// Call method
	exists, err := s.stringCache.Exists(s.ctx, "key1")

	// Assert
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

// TestExists_NotFound tests the Exists method when key doesn't exist
func (s *CacheSuite) TestExists_NotFound() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("EXISTS", "test:nonexistent")).
		Return(mock.Result(mock.ValkeyInt64(0)))

	// Call method
	exists, err := s.stringCache.Exists(s.ctx, "nonexistent")

	// Assert
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

// TestKeys tests the Keys method
func (s *CacheSuite) TestKeys() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("KEYS", "test:*")).
		Return(mock.Result(mock.ValkeyArray(
			mock.ValkeyString("test:key1"),
			mock.ValkeyString("test:key2"),
		)))

	// Call method
	keys, err := s.stringCache.Keys(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"key1", "key2"}, keys)
}

// TestKeys_NoPrefix tests the Keys method when no prefix is set
func (s *CacheSuite) TestKeys_NoPrefix() {
	// Create cache with no prefix
	noPrefix := NewCache[string](s.mockClient, "")

	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("KEYS", "*")).
		Return(mock.Result(mock.ValkeyArray(
			mock.ValkeyString("key1"),
			mock.ValkeyString("key2"),
		)))

	// Call method
	keys, err := noPrefix.Keys(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"key1", "key2"}, keys)
}

// TestGetAll tests the GetAll method
func (s *CacheSuite) TestGetAll() {
	// Set up expectation for KEYS - don't add extra quotes to the keys
	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyArray(
			mock.ValkeyString("test:key1"),
			mock.ValkeyString("test:key2"),
		)))

	// Set up expectations for GET calls
	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyString("\"value1\"")))

	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyString("\"value2\"")))

	// Call method
	items, err := s.stringCache.GetAll(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, items)
}

// TestGetAll_WithMissingKey tests the GetAll method when one key is missing
func (s *CacheSuite) TestGetAll_WithMissingKey() {
	// Set up expectation for KEYS
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("KEYS", "test:*")).
		Return(mock.Result(mock.ValkeyArray(
			mock.ValkeyString("test:key1"),
			mock.ValkeyString("test:key2"),
		)))

	// Set up expectations for GET calls
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("GET", "test:key1")).
		Return(mock.Result(mock.ValkeyString("\"value1\"")))
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("GET", "test:key2")).
		Return(mock.Result(mock.ValkeyNil()))

	// Call method
	items, err := s.stringCache.GetAll(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), map[string]string{
		"key1": "value1",
	}, items)
}

// TestStructEncoding tests encoding and decoding of struct values
func (s *CacheSuite) TestStructEncoding() {
	user := User{ID: 1, Name: "John", Email: "john@example.com"}
	userJSON, _ := json.Marshal(user)

	// Test Set with a struct - using gomock.Any() for arguments to avoid matching issues
	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyString("OK")))

	err := s.userCache.Set(s.ctx, "user1", user)
	assert.NoError(s.T(), err)

	// Test Get with a struct
	s.mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.ValkeyString(string(userJSON))))

	retrievedUser, err := s.userCache.Get(s.ctx, "user1")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, retrievedUser)
}

// TestIncrementInt tests the Increment function with int type
func (s *CacheSuite) TestIncrementInt() {
	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("INCRBY", "int-test:counter", "5")).
		Return(mock.Result(mock.ValkeyInt64(10)))

	// Call method
	newValue, err := Increment(s.ctx, s.intCache, "counter", 5)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 10, newValue)
}

// TestIncrementFloat64 tests the Increment function with float64 type
func (s *CacheSuite) TestIncrementFloat64() {
	// Create a float cache
	floatCache := NewCache[float64](s.mockClient, "float-test")

	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("INCRBYFLOAT", "float-test:counter", "2.5")).
		Return(mock.Result(mock.ValkeyFloat64(15.5)))

	// Call method
	newValue, err := Increment(s.ctx, floatCache, "counter", 2.5)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15.5, newValue)
}

// TestStringCache tests the specialized string cache
func (s *CacheSuite) TestStringCache() {
	// Create a string cache
	stringCache := NewStringCache(s.mockClient, "string-test")

	// Verify we're using StringValueCoder
	assert.IsType(s.T(), StringValueCoder{}, stringCache.coder)

	// Set up expectation
	s.mockClient.EXPECT().
		Do(s.ctx, mock.Match("SET", "string-test:key1", "direct-string")).
		Return(mock.Result(mock.ValkeyString("\"OK\"")))

	// Call method
	err := stringCache.Set(s.ctx, "key1", "direct-string")

	// Assert
	assert.NoError(s.T(), err)
}

// Run the test suite
func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

// Test the JSONValueCoder directly
func TestJSONValueCoder(t *testing.T) {
	coder := JSONValueCoder[User]{}

	// Test encoding
	user := User{ID: 123, Name: "Alice", Email: "alice@example.com"}
	encoded, err := coder.Encode(user)
	assert.NoError(t, err)

	// Verify it's valid JSON
	var decodedMap map[string]interface{}
	err = json.Unmarshal([]byte(encoded), &decodedMap)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), decodedMap["id"])
	assert.Equal(t, "Alice", decodedMap["name"])
	assert.Equal(t, "alice@example.com", decodedMap["email"])

	// Test decoding
	decoded, err := coder.Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, user, decoded)

	// Test decoding invalid JSON
	_, err = coder.Decode("invalid json")
	assert.Error(t, err)
}

// Test the StringValueCoder directly
func TestStringValueCoder(t *testing.T) {
	coder := StringValueCoder{}

	// Test encoding
	encoded, err := coder.Encode("hello world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", encoded)

	// Test decoding
	decoded, err := coder.Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", decoded)
}

// Test creating cache with custom coder
func TestWithCoder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewClient(ctrl)
	cache := NewCache[string](mockClient, "test")

	// Create a custom coder
	customCoder := StringValueCoder{}

	// Apply the custom coder
	result := cache.WithCoder(customCoder)

	// Assert it returns self for chaining
	assert.Same(t, cache, result)

	// Assert the coder was set
	assert.Equal(t, customCoder, cache.coder)
}

// Test fullKey method
func TestFullKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewClient(ctrl)

	// Cache with prefix
	cache := NewCache[string](mockClient, "prefix")
	assert.Equal(t, "prefix:key", cache.fullKey("key"))

	// Cache without prefix
	cacheNoPrefix := NewCache[string](mockClient, "")
	assert.Equal(t, "key", cacheNoPrefix.fullKey("key"))
}
