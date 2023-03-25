package cache

import "errors"

type Cacher[K comparable, T any] interface {
	// Get returns the value associated with the key.
	// Returns true if an eviction occurred.
	Get(key K) (T, bool)
	// Set sets the value associated with the key.
	// Returns true if the value was set.
	Set(key K, val T) bool
	// Delete deletes the value associated with the key.
	Delete(key K)
	// IsExist returns true if the key exists.
	IsExist(key K) bool
	// ClearAll clears all cache.
	ClearAll()

	// GetMulti returns the values associated with the keys.
	GetMulti(keys []K) map[K]T
	//PeekRandom returns a random value.
	PeekRandom() (T, error)
	// PeekRandomMulti returns random values.
	PeekRandomMulti(count int) ([]T, error)

	// SetMulti sets the values associated with the keys.
	SetMulti(kvs map[K]T) []bool
	// DeleteMulti deletes the values associated with the keys.
	DeleteMulti(keys []K)

	// ShouldCache returns true if the value should be cached.
	//ShouldCache(key K, val T) bool

	Len() int

	Cap() int

	// Resize resizes the cache.
	// Resize(cap int)
}

var (
	// 缓存为空
	ErrorCacheEmpty = errors.New("cache is empty")
)
