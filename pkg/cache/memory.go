package cache

import (
	"sync"
	"time"
)

type entry struct {
	value      interface{}
	expiration time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
	max   int
}

func NewMemoryCache(ttl time.Duration, max int) *MemoryCache {
	if max <= 0 {
		max = 1000
	}
	return &MemoryCache{
		items: make(map[string]entry, max),
		ttl:   ttl,
		max:   max,
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiration) {
		return nil, false
	}
	return e.value, true
}

func (c *MemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.items) >= c.max {
		//  clear all expired held in cache if reserved space is full
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expiration) {
				delete(c.items, k)
			}
		}
		// if still full, skip insert
		if len(c.items) >= c.max {
			return
		}
	}
	c.items[key] = entry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}
