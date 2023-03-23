package cache

import "errors"

type Cacher[K comparable, T any] interface {
	// Get returns the value associated with the key.
	Get(key K) (T, bool)
	// Set sets the value associated with the key.
	Set(key K, val T) bool
	// Delete deletes the value associated with the key.
	Delete(key K)
	// IsExist returns true if the key exists.
	IsExist(key K) bool
	// ClearAll clears all cache.
	ClearAll()

	// GetMulti returns the values associated with the keys.
	GetMulti(keys []K) map[K]T
	//GetRandom returns a random value.
	GetRandom() (T, error)
	// GetRandomMulti returns random values.
	GetRandomMulti(count int) ([]T, error)

	// SetMulti sets the values associated with the keys.
	SetMulti(kvs map[K]T) []bool
	// DeleteMulti deletes the values associated with the keys.
	DeleteMulti(keys []K)

	// ShouldCache returns true if the value should be cached.
	//ShouldCache(key K, val T) bool

	Len() int

	Cap() int

	// Resize resizes the cache.
	// Resize(cap int64)
}

var (
	// 缓存为空
	ErrorCacheEmpty = errors.New("cache is empty")
)
