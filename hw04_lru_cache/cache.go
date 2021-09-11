package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   string
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     *sync.RWMutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	item, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(item)
		item.Value.(*cacheItem).value = value
		return true
	}

	if l.queue.Len() == l.capacity {
		back := l.queue.Back()
		l.queue.Remove(back)
		delete(l.items, Key(back.Value.(*cacheItem).key))
	}

	item = l.queue.PushFront(&cacheItem{
		key:   string(key),
		value: value,
	})
	l.items[key] = item
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	item, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(item)
	return item.Value.(*cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		lock:     &sync.RWMutex{},
	}
}
