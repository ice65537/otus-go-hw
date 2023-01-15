package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
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

func (x *lruCache) Set(key Key, value interface{}) bool {
	x.Lock()
	defer x.Unlock()
	item, found := x.items[key]
	if found {
		x.queue.MoveToFront(item)
		(*item).Value = cacheItem{key: key, value: value}
	} else {
		x.queue.PushFront(cacheItem{key: key, value: value})
		x.items[key] = x.queue.Front()
		//
		if x.capacity < x.queue.Len() {
			ci := (*x.queue.Back()).Value.(cacheItem)
			delete(x.items, ci.key)
			x.queue.Remove(x.queue.Back())
		}
	}
	return found
}

func (x *lruCache) Get(key Key) (interface{}, bool) {
	x.Lock()
	defer x.Unlock()
	item, found := x.items[key]
	if found {
		x.queue.MoveToFront(item)
		ci := (*item).Value.(cacheItem)
		return ci.value, true
	}
	return nil, false
}

func (x *lruCache) Clear() {
	x.Lock()
	defer x.Unlock()
	for key, item := range x.items {
		x.queue.Remove(item)
		delete(x.items, key)
	}
}
