package shared

import (
	"sync"
)

type Cache struct {
	cache    map[string]interface{}
	requests chan request
}

type request struct {
	key      string
	value    interface{}
	response chan interface{}
}

var ccCacheInstanceLock sync.Once
var ccCacheInstance *Cache

func StaticCache() *Cache {
	ccCacheInstanceLock.Do(func() {
		ccCacheInstance = NewCache()
	})

	return ccCacheInstance
}

func NewCache() *Cache {
	cache := &Cache{requests: make(chan request)}
	go cache.serve()

	return cache
}

func (c *Cache) GetOrSet(key string, value interface{}) interface{} {
	resChan := make(chan interface{})
	c.requests <- request{key, value, resChan}
	response := <-resChan
	close(resChan)

	return response
}

func (c *Cache) Terminate() {
	close(c.requests)
	c.cache = nil
}

func (c *Cache) serve() {
	c.cache = make(map[string]interface{})

	for request := range c.requests {
		if request.value != nil {
			c.cache[request.key] = request.value
			request.response <- request.value
		} else {
			request.response <- c.cache[request.key]
		}
	}
}
