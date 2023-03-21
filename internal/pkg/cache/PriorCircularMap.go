package cache

import (
	"runtime"
	"sync"

	"github.com/Doraemonkeys/douyin2/pkg/third_party/priorityQueue"
)

// 获取权重的接口
type Weighter[T any] interface {
	Weight() int64
	Less(T) bool
}

type pairs[K comparable, T Weighter[T]] struct {
	Key K
	Val T
}

func (p pairs[K, T]) Less(other pairs[K, T]) bool {
	return p.Val.Less(other.Val)
}

func newPairs[K comparable, T Weighter[T]](key K, val T) pairs[K, T] {
	return pairs[K, T]{Key: key, Val: val}
}

type PriorCircularMap[K comparable, T Weighter[T]] struct {
	cache    map[K]*T
	capacity int
	//priorityQueu is not thread safe, so we also need a lock ,so share the same lock with cache
	cacheLock sync.RWMutex
	// plan to persist
	planPersist     map[K]*T
	persistLock     sync.Mutex
	persistenter    Persistenter[K, T]
	priorityQueue   *priorityQueue.PriorityQueue[pairs[K, T]]
	minWeightVal    *T
	minWeightValKey K
}

func NewPriorCircularMap[K comparable, T Weighter[T]](cap int) *PriorCircularMap[K, T] {
	return &PriorCircularMap[K, T]{
		cache: make(map[K]*T, cap),
		//cap/10+1 is magic number, it is the best value i guess
		planPersist:   make(map[K]*T, cap/10+1),
		cacheLock:     sync.RWMutex{},
		persistLock:   sync.Mutex{},
		capacity:      cap,
		priorityQueue: priorityQueue.NewPriorityQueue[pairs[K, T]](),
	}
}

// SetPersistent sets the persistent of the cache.
func (c *PriorCircularMap[K, T]) SetPersistent(persistent Persistenter[K, T]) {
	c.persistenter = persistent
}

func (c *PriorCircularMap[K, T]) Cap() int {
	return c.capacity
}

func (c *PriorCircularMap[K, T]) Len() int {
	return len(c.cache)
}

// Get returns the value associated with the key.
func (c *PriorCircularMap[K, T]) Get(key K) (T, bool) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	v, ok := c.cache[key]
	if !ok {
		var val T
		return val, false
	}
	return *v, true
}

// Set sets the value associated with the key.
// If the key is already in the cache, it will be replaced.
// If the cache is full, only the value's weight is greater than the minimum weight in the cache,
// the value will be replaced.Otherwise, the value will not be set, and the function will return false.
func (c *PriorCircularMap[K, T]) Set(key K, val T) bool {
	if c.Len() < c.Cap() {
		c.cacheLock.Lock()
		c.cache[key] = &val
		c.priorityQueue.Push(newPairs(key, val))
		if c.minWeightVal == nil || val.Less((*c.minWeightVal)) {
			c.minWeightVal = &val
			c.minWeightValKey = key
		}
		c.cacheLock.Unlock()

		c.persistLock.Lock()
		c.planPersist[key] = &val
		c.persistLock.Unlock()
		return true
	}
	// cache is full
	if val.Less((*c.minWeightVal)) {
		return false
	}
	c.cacheLock.Lock()
	delete(c.cache, c.minWeightValKey)
	c.cache[key] = &val
	c.priorityQueue.Pop()
	c.priorityQueue.Push(newPairs(key, val))
	var temp = c.priorityQueue.Top()
	c.minWeightVal = &temp.Val
	c.minWeightValKey = temp.Key
	c.cacheLock.Unlock()

	c.persistLock.Lock()
	c.planPersist[key] = &val
	c.persistLock.Unlock()
	return true
}

// Delete deletes the value associated with the key.
func (c *PriorCircularMap[K, T]) Delete(key K) {
	c.cacheLock.Lock()
	delete(c.cache, key)
	c.cacheLock.Unlock()
}

// IsExist returns true if the key exists.
func (c *PriorCircularMap[K, T]) IsExist(key K) bool {
	c.cacheLock.RLock()
	_, ok := c.cache[key]
	c.cacheLock.RUnlock()
	return ok
}

// ClearAll clears all cache.
func (c *PriorCircularMap[K, T]) ClearAll() {
	c.cacheLock.Lock()
	c.cache = nil
	runtime.GC()
	c.priorityQueue = priorityQueue.NewPriorityQueue[pairs[K, T]]()
	c.cache = make(map[K]*T, c.capacity)
	c.cacheLock.Unlock()
}

// GetMult returns the values associated with the keys.
// If the key is not in the cache, it will not return.
func (c *PriorCircularMap[K, T]) GetMulti(keys []K) map[K]T {
	var result = make(map[K]T, len(keys))
	c.cacheLock.RLock()
	for _, key := range keys {
		if val, ok := c.cache[key]; ok {
			result[key] = *val
		}
	}
	c.cacheLock.RUnlock()
	return result
}

// // SetMulti sets the values associated with the keys.
// func (c *CirCurMap) SetMulti(keyvals map[cache.K]cache.T) error {
// 	panic("not implemented") // TODO: Implement
// }

// // DeleteMulti deletes the values associated with the keys.
// func (c *CirCurMap) DeleteMulti(keys []cache.K) error {
// 	panic("not implemented") // TODO: Implement
// }

// // SetExpire sets the value associated with the key and expire time.
// func (c *CirCurMap) SetExpire(key cache.K, val cache.T, timeout int64) error {
// 	panic("not implemented") // TODO: Implement
// }

// // GetCapacity returns the capacity of the cache.
// func (c *CirCurMap) GetCapacity() int64 {
// 	panic("not implemented") // TODO: Implement
// }

// // GetUsed returns the used of the cache.
// func (c *CirCurMap) GetUsed() int64 {
// 	panic("not implemented") // TODO: Implement
// }

// // EnablePersistent enables the persistent of the cache.
// func (c *CirCurMap) EnablePersistent(enable bool) error {
// 	panic("not implemented") // TODO: Implement
// }

// // IsPersistentEnabled returns true if the persistent of the cache is enabled.
// func (c *CirCurMap) IsPersistentEnabled() bool {
// 	panic("not implemented") // TODO: Implement
// }

// // GetPersistent returns the persistent of the cache.
// func (c *CirCurMap) GetPersistent() cache.Persistenter[cache.K, cache.T] {
// 	panic("not implemented") // TODO: Implement
// }

// // SetPersistCycle sets the persist cycle of the cache.
// func (c *CirCurMap) SetPersistCycle(cycle time.Duration) error {
// 	panic("not implemented") // TODO: Implement
// }
