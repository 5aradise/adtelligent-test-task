package cache

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Store(K, V)
	Load(K) (V, bool)
	Delete(K)
}

type cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
	l    *slog.Logger
}

func New[K comparable, V any](logger *slog.Logger, loadValuesFunc func([]K) (map[K]V, error), updateDelay time.Duration) *cache[K, V] {
	c := &cache[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V),
		l:    logger,
	}
	go c.run(loadValuesFunc, updateDelay)
	return c
}

func (c *cache[K, V]) run(loadValues func([]K) (map[K]V, error), updateDelay time.Duration) {
	t := time.NewTicker(updateDelay)
	var err error
	c.l.Info("start cache", slog.Duration("update delay", updateDelay))
	for range t.C {
		err = c.update(loadValues)
		if err != nil {
			c.l.Error("failed to update cache", slog.String("error", err.Error()))
		} else {
			c.l.Info("cache successfully updated")
		}
	}
}

func (c *cache[K, V]) update(loadValues func([]K) (map[K]V, error)) error {
	const op = "cache.update"

	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]K, len(c.data))
	var i int
	for k := range c.data {
		keys[i] = k
		i++
	}

	newData, err := loadValues(keys)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.data = newData
	c.l.Debug("cache data replaced", slog.Any("new data", newData))
	return nil
}

func (c *cache[K, V]) Store(key K, value V) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
	c.l.Debug("value stored in cache", slog.Any("key", key), slog.Any("value", value))
}

func (c *cache[K, V]) Load(key K) (value V, ok bool) {
	c.mu.RLock()
	value, ok = c.data[key]
	c.mu.RUnlock()
	if ok {
		c.l.Debug("value loaded from cache", slog.Any("key", key), slog.Any("value", value))
	}
	return
}

func (c *cache[K, V]) Delete(key K) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	c.l.Debug("value deleted from cache", slog.Any("key", key))
}
