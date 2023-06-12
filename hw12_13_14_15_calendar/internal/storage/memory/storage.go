package memorystorage

import "sync"

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func New(stgType string) *Storage {
	return &Storage{}
}

// TODO
