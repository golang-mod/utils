package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	c          *cache.Cache
	countLimit int
}

// NewCache cleanupInterval 定时清理缓存的时间
func NewCache(countLimit int, timeLimit, cleanupInterval time.Duration) *Cache {
	res := new(Cache)
	res.c = cache.New(timeLimit, cleanupInterval)
	res.countLimit = countLimit

	return res
}

func (c *Cache) Limit(key string) bool {
	count, ok := c.c.Get(key)
	if !ok {
		return false
	}
	curCount := count.(int)
	if curCount >= c.countLimit {
		return true
	}
	return false
}

func (c *Cache) Add(key string) {
	count, expiration, ok := c.c.GetWithExpiration(key)
	if !ok {
		c.c.SetDefault(key, 1)
		return
	}

	// 这里的过期时间重新算一下
	duration := expiration.Sub(time.Now())
	if duration <= 0 {
		return
	}

	curCount := count.(int)
	c.c.Set(key, curCount+1, duration)
	return
}
