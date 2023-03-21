package cache

import "time"

type Persistenter[K comparable, T any] interface {
	// Get returns the value associated with the key.
	Get(key K) (T, error)
	// Set sets the value associated with the key.
	Set(key K, val T) error
	// Delete deletes the value associated with the key.
	Delete(key K) error
	// IsExist returns true if the key exists.
	IsExist(key K) bool
	// ClearAll clears all cache.
	ClearAll() error
	// StartAndGC starts the cache gc process.
	//StartAndGC(config *Config) error
	// GetMulti returns the values associated with the keys.
	GetMulti(keys []K) ([]T, error)
	// SetMulti sets the values associated with the keys.
	SetMulti(keyvals map[K]T) error
	// DeleteMulti deletes the values associated with the keys.
	DeleteMulti(keys []K) error
}

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

	// StartAndGC starts the cache gc process.
	//StartAndGC(config *Config) error

	// GetMulti returns the values associated with the keys.
	GetMulti(keys []K) map[K]T
	//GetRandom returns a random value.
	GetRandom() (T, error)
	// GetRandomMulti returns random values.
	GetRandomMulti(count int) ([]T, error)

	// SetMulti sets the values associated with the keys.
	SetMulti(keyvals []T) error
	// DeleteMulti deletes the values associated with the keys.
	DeleteMulti(keys []K) error
	// SetExpire sets the value associated with the key and expire time.
	SetExpire(key K, val T, timeout int64) error

	// ShouldCache returns true if the value should be cached.
	ShouldCache(key K, val T) bool

	// Incr increases counter associated with the key by delta.
	//Incr(key K, delta int64) error
	// Decr decreases counter associated with the key by delta.
	//Decr(key K, delta int64) error
	// IsTimeout returns true if err is timeout.
	//IsTimeout(err error) bool
	// SetCapacity sets the capacity of the cache.
	//SetCapacity(capacity int64) error

	// GetCapacity returns the capacity of the cache.
	GetCapacity() int64
	// GetUsed returns the used of the cache.
	GetUsed() int64
	// EnablePersistent enables the persistent of the cache.
	EnablePersistent(enable bool) error
	// IsPersistentEnabled returns true if the persistent of the cache is enabled.
	IsPersistentEnabled() bool
	// GetPersistent returns the persistent of the cache.
	GetPersistent() Persistenter[K, T]
	// SetPersistent sets the persistent of the cache.
	SetPersistent(persistent Persistenter[K, T])
	// SetPersistCycle sets the persist cycle of the cache.
	SetPersistCycle(cycle time.Duration) error
}
