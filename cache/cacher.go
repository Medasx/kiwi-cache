package cache

import (
	"kiwi-entry-task/fetcher"
	"sync"
	"time"
)

type Cacher interface {
	Get(id int) (string, error)
}

// itemsMap represents concurrent safe map for our purposes `map[int]string`
type itemsMap struct {
	sync.RWMutex
	data map[int]string
}

// load string value from map
func (m *itemsMap) load(key int) (string, bool) {
	m.RLock()
	defer m.RUnlock()
	if value, ok := m.data[key]; ok {
		return value, true
	}
	return "", false
}

// store into map
func (m *itemsMap) store(key int, value string) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = value
}

// replaceItems replaces all data with new ones
func (m *itemsMap) replaceItems(bulk map[int]string) {
	m.Lock()
	defer m.Unlock()
	m.data = bulk
}

// cache implements Cacher interface
type cache struct {
	fetcher fetcher.Fetcher
	items   *itemsMap
}

// NewCache initializes new cache
func NewCache(fetcher fetcher.Fetcher, refreshInterval time.Duration, errorHandler func(error)) Cacher {
	c := cache{
		fetcher: fetcher,
		items:   &itemsMap{data: make(map[int]string)},
	}
	// fetch and store all data
	if err := c.reloadData(); err != nil {
		errorHandler(err)
	}
	go func() {
		// reload data after defined time period
		for range time.After(refreshInterval) {
			if err := c.reloadData(); err != nil {
				errorHandler(err)
			}
		}
	}()
	return c
}

// Get returns a value of given key loaded locally if exists, otherwise requested and stored from fetcher
func (c cache) Get(key int) (string, error) {
	if val, ok := c.items.load(key); ok {
		return val, nil
	}

	val, err := c.fetcher.Fetch(key)
	if err != nil {
		return "", err
	}

	c.items.store(key, val)
	return val, nil
}

// reloadData replace stored data with currently fetched data
func (c cache) reloadData() error {
	data, err := c.fetcher.FetchAll()
	if err != nil {
		return err
	}

	c.items.replaceItems(data)
	return nil
}
