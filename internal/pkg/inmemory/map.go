package inmemory

import (
	"sync"
)

type InMemory[keyType comparable, valueType any] interface {
	GetAll() map[keyType]valueType
	RewriteAll(in map[keyType]valueType)
	AddData(key keyType, value valueType)
	Get(key keyType) (value valueType, ok bool)
	Len() int
	DeleteKey(key keyType)
}

type inMemory[k comparable, v any] struct {
	mutex sync.Mutex
	data  map[k]v
}

func New[keyType comparable, valueType any]() InMemory[keyType, valueType] {
	return &inMemory[keyType, valueType]{
		data: make(map[keyType]valueType),
	}
}

func (c *inMemory[k, v]) RewriteAll(in map[k]v) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if in == nil {
		in = make(map[k]v)
	}
	c.data = in
}

func (c *inMemory[k, v]) AddData(key k, value v) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *inMemory[k, v]) DeleteKey(key k) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

func (c *inMemory[k, v]) Get(key k) (value v, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok = c.data[key]

	return value, ok
}

func (c *inMemory[k, v]) GetAll() map[k]v {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outMap := make(map[k]v, len(c.data))
	for key, val := range c.data {
		outMap[key] = val
	}
	return outMap
}

func (c *inMemory[k, v]) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.data)
}
