package cache

import (
	"container/list"
	"sync"
	"time"
)

/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package lru implements an LRU c.

// c is an LRU c. It is not safe for concurrent access.
type cache struct {
	sync.RWMutex

	// MaxEntries is the maximum number of c entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int

	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the c.
	OnEvicted func(key, value any)

	ll    *list.List
	cache map[any]*list.Element

	evictedCache *sync.Map

	evictedRate int
}

type entry struct {
	key   any
	value any
}

type Option func(c *cache)

func SetEvictedRate(rate int) Option {
	return func(c *cache) {
		c.evictedRate = rate
	}
}

// New creates a new c.
// If maxEntries is zero, the c has no limit and it's assumed
// that eviction is done by the caller.
func New(maxEntries int, opts ...Option) *cache {
	c := &cache{
		MaxEntries:   maxEntries,
		ll:           list.New(),
		cache:        make(map[any]*list.Element),
		evictedCache: &sync.Map{},
		evictedRate:  maxEntries / 1000,
	}
	if c.evictedRate < 5000 {
		c.evictedRate = 5000
	}
	for _, opt := range opts {
		opt(c)
	}
	c.evicted()
	return c
}

func (c *cache) evicted() {
	go func() {
		limiter := make(chan any, c.evictedRate)
		for {
			c.evictedCache.Range(func(key, value any) bool {
				limiter <- nil
				go func() {
					defer func() {
						<-limiter
						if e := recover(); e != nil {

						}
					}()
					c.evictedCache.Delete(key)
					if c.OnEvicted != nil {
						c.OnEvicted(key, value)
					}
				}()
				return true
			})
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Add adds a value to the c.
func (c *cache) Add(key, value any) {

	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		c.cache = make(map[any]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.removeOldest()
	}
}

// Get looks up a key's value from the c.
func (c *cache) Get(key any) (value any, ok bool) {

	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the c.
func (c *cache) Remove(key any) any {

	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		return nil
	}
	if ele, hit := c.cache[key]; hit {
		return c.removeElement(ele)
	}
	return nil
}

// RemoveOldest removes the oldest item from the c.
func (c *cache) removeOldest() any {

	if c.cache == nil {
		return nil
	}
	ele := c.ll.Back()

	if nil != ele {
		return c.removeElement(ele)
	}
	return nil
}

func (c *cache) removeElement(e *list.Element) any {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		// c.OnEvicted(kv.key, kv.value)
		c.evictedCache.LoadOrStore(kv.key, kv.value)
	}
	return kv.value
}

// Len returns the number of items in the c.
func (c *cache) Len() int {
	c.RLock()
	defer c.RUnlock()
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the c.
func (c *cache) Clear() {
	c.Lock()
	defer c.Unlock()

	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			// c.OnEvicted(kv.key, kv.value)
			c.evictedCache.LoadOrStore(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}

// iterator
func (c *cache) Iterator(do func(k, v any) error) {
	c.RLock()
	defer c.RUnlock()
	for _, e := range c.cache {
		kv := e.Value.(*entry)
		err := do(kv.key, kv.value)
		//err break finish iterator
		if nil != err {
			break
		}
	}
}
