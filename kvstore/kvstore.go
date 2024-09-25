package kvstore

import (
	"fmt"
	"sync"
)

// KeyValueStore defines the interface for a key-value store.
type KeyValueStore interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) bool
}

// InMemoryKVStore implements the KeyValueStore interface using sync.Map.
type InMemoryKVStore struct {
	store sync.Map
}

// Set sets a value for a given key. This is kept for backward compatibility and can be removed if not needed.
func (kv *InMemoryKVStore) Set(key string, value interface{}) error {
	kv.store.Store(key, value)
	return nil
}

// Get retrieves the string value for a given key.
func (kv *InMemoryKVStore) Get(key string) (string, error) {
	value, found := kv.store.Load(key)
	if !found {
		return "", fmt.Errorf("key not found")
	}
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}
	return strValue, nil
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
