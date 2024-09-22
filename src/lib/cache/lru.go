package cache

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] struct {
	l   *list.List
	m   map[K]*list.Element
	mu  sync.Mutex
	cap int
}

type node[K comparable] struct {
	key   K
	value any
}

func NewLRU[K comparable, V any](cap int) *LRU[K, V] {
	return &LRU[K, V]{
		cap: cap,
		l:   &list.List{},
		m:   make(map[K]*list.Element),
	}
}

func (this *LRU[K, V]) Set(key K, value V) {
	this.mu.Lock()
	defer this.mu.Unlock()

	elem, has := this.m[key]
	if has {
		this.l.Remove(elem)
	}

	new := this.l.PushFront(node[K]{
		key:   key,
		value: value,
	})
	this.m[key] = new
  
	len := this.l.Len()
	if len > this.cap {
		back := this.l.Back()
		this.l.Remove(back)
    node := back.Value.(node[K])
    delete(this.m, node.key)
	}
}

func (this *LRU[K, T]) Get(key K) (T, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()

	elem, has := this.m[key]
	if !has {
    var zero T
		return zero, false
	}

	this.l.MoveToFront(elem)
  node := elem.Value.(node[K])
  value := node.value.(T)
	return value, true
}

func (this *LRU[K, T]) GetAll(keys []K) []T {
  out := make([]T, len(keys))
  
  for i, key := range keys {
    value, has := this.Get(key)
    if !has {
      return []T{}
    }

    out[i] = value
  } 

  return out
}
