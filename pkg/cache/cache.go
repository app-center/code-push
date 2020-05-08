package cache

import (
	"container/list"
	"sync"
)

type KeyType = interface{}
type ValueType = interface{}
type ValueAllocationFunc func(key KeyType) (value ValueType, ok bool)

var defaultValueAllocationFunc = func(key KeyType) (value ValueType, ok bool) {
	ok = false
	return
}

type entry struct {
	Key   KeyType
	Value ValueType
}

type Cache struct {
	capacity  int
	allocFunc ValueAllocationFunc
	mu        sync.RWMutex
	list      *list.List
	elmMap    map[KeyType]*list.Element
}

// 查询缓存中是否存在 key，
// 存在则 ok = true，并返回 value
// 不存在则 ok = false
func (c *Cache) Query(key KeyType) (value ValueType, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elm, ok := c.elmMap[key]; ok {
		c.list.MoveToFront(elm)
		value = elm.Value.(*entry).Value
		return value, true
	}

	return nil, false
}

// 查询缓存中是否存在 key，
// 存在则 ok = true，并返回 value
// 不存在，则尝试调用初始化函数生成 value 并存入缓存中，存入成功，也会返回 true
func (c *Cache) Get(key KeyType) (ValueType, bool) {
	c.mu.RLock()

	if elm, ok := c.elmMap[key]; ok {
		c.list.MoveToFront(elm)
		c.mu.RUnlock()
		return elm.Value.(*entry).Value, true
	}

	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.allocFunc == nil {
		return nil, false
	}

	newValue, newAllocated := c.allocFunc(key)
	if !newAllocated {
		return nil, false
	}

	c.set(key, newValue)

	return newValue, true
}

func (c *Cache) Set(key KeyType, value ValueType) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	evictedEntry := c.set(key, value)

	return evictedEntry != nil
}

func (c *Cache) Contains(key KeyType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.contains(key)
}

func (c *Cache) Remove(key KeyType) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elm, ok := c.elmMap[key]; elm != nil && ok {
		c.removeElement(elm)
		return true
	}

	return false
}

func (c *Cache) Capacity() int {
	return c.capacity
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.len()
}

func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.elmMap {
		delete(c.elmMap, key)
	}

	c.list.Init()
}

func (c *Cache) set(key KeyType, value ValueType) (evictEntry *entry) {
	if c.contains(key) {
		elm := c.elmMap[key]
		elm.Value.(*entry).Value = value
		c.list.MoveToFront(elm)
		return
	}

	elm := c.list.PushFront(&entry{
		Key:   key,
		Value: value,
	})

	c.elmMap[key] = elm

	if c.len() > c.capacity {
		oldestElm := c.removeOldest()
		if oldestElm != nil {
			evictEntry = oldestElm.Value.(*entry)
		}
	}

	return
}

func (c *Cache) contains(key KeyType) bool {
	_, ok := c.elmMap[key]
	return ok
}

func (c *Cache) len() int {
	return c.list.Len()
}

func (c *Cache) removeOldest() *list.Element {
	if elm := c.list.Back(); elm != nil {
		c.removeElement(elm)
		return elm
	}

	return nil
}

func (c *Cache) removeElement(elm *list.Element) {
	c.list.Remove(elm)
	kv := elm.Value.(*entry)

	delete(c.elmMap, kv.Key)
}

type CtorConfig struct {
	Capacity  int
	AllocFunc ValueAllocationFunc
}

func New(config CtorConfig) *Cache {
	if config.AllocFunc == nil {
		config.AllocFunc = defaultValueAllocationFunc
	}

	if config.Capacity < 0 {
		config.Capacity = 0
	}

	return &Cache{
		capacity:  config.Capacity,
		allocFunc: config.AllocFunc,
		list:      list.New(),
		elmMap:    make(map[KeyType]*list.Element),
	}
}
