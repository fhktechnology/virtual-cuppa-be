package utils

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	Expiration time.Time
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

var confirmCodeCache *Cache

func init() {
	confirmCodeCache = &Cache{
		items: make(map[string]CacheItem),
	}
	go confirmCodeCache.cleanupExpired()
}

func GetConfirmCodeCache() *Cache {
	return confirmCodeCache
}

func (c *Cache) Set(key string, value string, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists {
		return "", false
	}
	if time.Now().After(item.Expiration) {
		return "", false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
