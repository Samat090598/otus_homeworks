package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.RWMutex
}

type pair struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, found := c.items[key]; found {
		v, _ := val.Value.(pair)
		v.value = value
		val.Value = v
		c.queue.MoveToFront(val)
		c.items[key] = val
		return true
	}

	if c.queue.Len() >= c.capacity {
		lastElem := c.queue.Back()
		v, _ := lastElem.Value.(pair)
		delete(c.items, v.key)
		c.queue.Remove(lastElem)
	}

	val := pair{key: key, value: value}
	item := c.queue.PushFront(val)
	c.items[key] = item

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, found := c.items[key]; found {
		c.queue.MoveToFront(val)
		v, _ := val.Value.(pair)
		return v.value, found
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue.Init()
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		delete(c.items, key)
	}
}
