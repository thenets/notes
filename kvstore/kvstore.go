package kvstore

import (
	"fmt"
	"sync"
)

// KeyValueStore defines the interface for a key-value store.
type KeyValueStore interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) bool
}

// InMemoryKVStore implements the KeyValueStore interface using sync.Map.
type InMemoryKVStore struct {
	store sync.Map
}

// Set sets a value for a given key.
func (kv *InMemoryKVStore) Set(key string, value interface{}) error {
	kv.store.Store(key, value)
	return nil
}

// Get retrieves the value for a given key.
func (kv *InMemoryKVStore) Get(key string) (interface{}, error) {
	value, found := kv.store.Load(key)
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	return value, nil
}

// Delete deletes a key-value pair.
func (kv *InMemoryKVStore) Delete(key string) error {
	kv.store.Delete(key)
	return nil
}

// Exists checks if a key exists in the store.
func (kv *InMemoryKVStore) Exists(key string) bool {
	_, found := kv.store.Load(key)
	return found
}
