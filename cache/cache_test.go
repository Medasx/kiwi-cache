package cache

import (
	"kiwi-entry-task/fetcher"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	f := fetcher.NewTestFetcher()
	c := NewCache(f, time.Minute, func(err error) {
		return
	})

	if val, err := c.Get(1); err != nil {
		t.Error(err)
	} else if val != "1" {
		t.Errorf("expected '1' got '%s'", val)
	}

	if val, err := c.Get(2); err != nil {
		t.Error(err)
	} else if val != "2" {
		t.Errorf("expected '2' got '%s'", val)
	}

	if val, err := c.Get(42); err != nil {
		t.Error(err)
	} else if val != "ultimate answer" {
		t.Errorf("expected 'ultimate answer' got '%s'", val)
	}
}
