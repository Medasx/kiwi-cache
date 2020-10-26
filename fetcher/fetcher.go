package fetcher

import "fmt"

type Fetcher interface {
    Fetch(id int) (string, error)
    FetchAll() (map[int]string, error)
}

func NewTestFetcher() Fetcher {
    return &testFetcher{data: map[int]string{
        1: "1",
        2: "2",
        42: "ultimate answer",
    }}
}

type testFetcher struct {
    data map[int]string
}

func (t testFetcher) Fetch(id int) (string, error) {
    if val, ok := t.data[id]; ok {
        return val, nil
    }
    return "", fmt.Errorf("not found")
}

func (t testFetcher) FetchAll() (map[int]string, error) {
    return t.data, nil
}

