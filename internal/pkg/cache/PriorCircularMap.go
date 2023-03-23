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
	Val *T
}

func (p pairs[K, T]) Less(other pairs[K, T]) bool {
	return (*p.Val).Less(*other.Val)
}

func newPairs[K comparable, T Weighter[T]](key K, val *T) pairs[K, T] {
	return pairs[K, T]{Key: key, Val: val}
}

// PriorCircularMap 通过Key索引，通过Value的权重排序，权重越大，缓存级别越高，越不容易被淘汰。
// 优先队列缓存存在一个问题，如果value频繁无规律变化，那么优先队列的排序会频繁变化，导致性能下降。
// 如果value的权重一直呈现上升或下降趋势，则在堆调整的过程中会比较迅速，影响不大。
type PriorCircularMap[K comparable, T Weighter[T]] struct {
	cache    map[K]*T
	capacity int
	//priorityQueu is not thread safe, so we also need a lock ,so share the same lock with cache
	cacheLock       sync.RWMutex
	priorityQueue   *priorityQueue.PriorityQueue[pairs[K, T]]
	minWeightVal    *T
	minWeightValKey K
}

func NewPriorCircularMap[K comparable, T Weighter[T]](cap int) *PriorCircularMap[K, T] {
	return &PriorCircularMap[K, T]{
		cache: make(map[K]*T, cap),
		//cap/10+1 is magic number, it is the best value i guess
		cacheLock:     sync.RWMutex{},
		capacity:      cap,
		priorityQueue: priorityQueue.NewPriorityQueue[pairs[K, T]](),
	}
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
	// cache is full
	if c.Len() == c.Cap() && val.Less((*c.minWeightVal)) {
		return false
	}
	_, exist := c.cache[key]
	if c.Len() < c.Cap() || exist {
		var valPtr = &val
		c.cacheLock.Lock()
		c.cache[key] = valPtr
		c.priorityQueue.Push(newPairs(key, valPtr))
		if c.minWeightVal == nil || val.Less((*c.minWeightVal)) {
			c.minWeightVal = valPtr
			c.minWeightValKey = key
		}
		c.cacheLock.Unlock()
		return true
	}

	// cache is full and val is greater than the minimum weight in the cache
	var valPtr = &val
	c.cacheLock.Lock()
	delete(c.cache, c.minWeightValKey)
	c.cache[key] = valPtr
	c.priorityQueue.Pop()
	c.priorityQueue.Push(newPairs(key, valPtr))
	var temp = c.priorityQueue.Top()
	c.minWeightVal = temp.Val
	c.minWeightValKey = temp.Key
	c.cacheLock.Unlock()

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

// SetMult sets the values associated with the keys.
// If the key is already in the cache, it will be replaced.
// If the cache is full, only the value's weight is greater than the minimum weight in the cache,
// the value will be replaced.Otherwise, the value will not be set, and the function will return false.
func (c *PriorCircularMap[K, T]) SetMulti(kvs map[K]T) []bool {
	var result = make([]bool, len(kvs))
	var i = 0
	for key, val := range kvs {
		result[i] = c.Set(key, val)
		i++
	}
	return result
}

// DeleteMult deletes the values associated with the keys.
func (c *PriorCircularMap[K, T]) DeleteMulti(keys []K) {
	c.cacheLock.Lock()
	for _, key := range keys {
		delete(c.cache, key)
	}
	c.cacheLock.Unlock()
}

// GetRandom returns a random value from the cache.
func (c *PriorCircularMap[K, T]) GetRandom() (T, error) {
	for _, v := range c.cache {
		return *v, nil
	}
	var zero T
	return zero, ErrorCacheEmpty
}

// GetRandomMulti returns a random values from the cache.
func (c *PriorCircularMap[K, T]) GetRandomMulti(count int) ([]T, error) {
	if count > c.Len() {
		count = c.Len()
	}
	var result = make([]T, count)
	var i = 0
	for _, v := range c.cache {
		result[i] = *v
		i++
		if i == count {
			break
		}
	}
	return result, nil
}
