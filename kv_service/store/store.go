package store

import (
	"sync"
)

// Store is a key-value store.
type Store interface {
	Set(key string, value any)
	Get(key string) any // returns nil if key not found
	Delete(key string)  // no-op if key does not exist
}

// inMemoryStore is a thread-safe, in-memory implementation
// of a Store.
type inMemoryStore struct {
	store map[string]any
	mu    sync.RWMutex
}

func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{store: make(map[string]any), mu: sync.RWMutex{}}
}

func (s *inMemoryStore) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

func (s *inMemoryStore) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.store[key]
	if !ok {
		return nil
	}
	return value
}

func (s *inMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
