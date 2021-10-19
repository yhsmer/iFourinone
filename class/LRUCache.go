package class

import (
	"container/list"
	"sync"
	"unsafe"
)

// Cache cache use LRU algorithm
type LRUCache struct {
	MaxBytes   int64 // 允许使用的最大内存
	NowBytes   int64 //当前已经使用的内存
	data       *list.List
	cache      map[string]*list.Element
	mu         sync.Mutex
}

type date struct {
	key   string
	value Value
}

type Value interface {
}

func NewLRUCache(maxBytes int64) *LRUCache {
	return &LRUCache{
		MaxBytes:  maxBytes,
		NowBytes:  0,
		data:      list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// Get get the key's value, move the node to the front of the queue, means its frequently visited
func (c *LRUCache) Get(key string) (value Value, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*date)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *LRUCache) removeOldest() {
	ele := c.data.Back()
	if ele != nil {
		c.data.Remove(ele)
		kv := ele.Value.(*date)
		delete(c.cache, kv.key)
		c.NowBytes -= int64(unsafe.Sizeof(kv.key)) + int64(unsafe.Sizeof(kv.value))
	}
}

// Add adds a value to the cache.
func (c *LRUCache) Add(key string, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if date exist, then replace
	if ele, ok := c.cache[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*date)
		c.NowBytes += int64(unsafe.Sizeof(value)) - int64(unsafe.Sizeof(kv.value))
		kv.value = value
	} else {
		ele := c.data.PushFront(&date{key, value})
		c.cache[key] = ele
		c.NowBytes += int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(value))
	}
	for c.MaxBytes != 0 && c.MaxBytes < c.NowBytes {
		c.removeOldest()
	}
}

// Len the number of cache entries
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.data.Len()
}
