package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
}

func (s *storeTestSuite) TestSet() {
	store := NewInMemoryStore()

	// Test setting a string value
	store.Set("key1", "value1")
	assert.Equal(s.T(), "value1", store.Get("key1"))

	// Test setting an integer value
	store.Set("key2", 42)
	assert.Equal(s.T(), 42, store.Get("key2"))

	// Test setting a boolean value
	store.Set("key3", true)
	assert.Equal(s.T(), true, store.Get("key3"))

	// Test setting a nil value
	store.Set("key4", nil)
	assert.Nil(s.T(), store.Get("key4"))

	// Test overwriting an existing key
	store.Set("key1", "new_value")
	assert.Equal(s.T(), "new_value", store.Get("key1"))
}

func (s *storeTestSuite) TestGet() {
	store := NewInMemoryStore()

	// Test getting a non-existent key
	assert.Nil(s.T(), store.Get("nonexistent"))

	// Test getting an existing key
	store.Set("existing", "value")
	assert.Equal(s.T(), "value", store.Get("existing"))

	// Test getting a key with nil value
	store.Set("nil_key", nil)
	assert.Nil(s.T(), store.Get("nil_key"))
}

func (s *storeTestSuite) TestDelete() {
	store := NewInMemoryStore()

	// Test deleting a non-existent key (should be no-op)
	store.Delete("nonexistent")
	assert.Nil(s.T(), store.Get("nonexistent"))

	// Test deleting an existing key
	store.Set("to_delete", "value")
	assert.Equal(s.T(), "value", store.Get("to_delete"))

	store.Delete("to_delete")
	assert.Nil(s.T(), store.Get("to_delete"))

	// Test deleting the same key again (should be no-op)
	store.Delete("to_delete")
	assert.Nil(s.T(), store.Get("to_delete"))
}

func (s *storeTestSuite) TestConcurrentAccess() {
	store := NewInMemoryStore()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d", id)
				value := fmt.Sprintf("value_%d", j)
				store.Set(key, value)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d", id)
				store.Get(key) // We don't check the value since it's changing concurrently
			}
		}(i)
	}

	// Concurrent deletes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d", id)
				store.Delete(key)
			}
		}(i)
	}

	wg.Wait()
	// If we reach here without panic or race condition, the test passes
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}
