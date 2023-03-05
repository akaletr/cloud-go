package storage

import (
	"errors"
	"sync"
)

type storage struct {
	sync.RWMutex
	data map[string]string
}

var ErrorNoSuchKey = errors.New("error: no such key")

func New() Storage {
	return &storage{
		RWMutex: sync.RWMutex{},
		data:    make(map[string]string),
	}
}

func (storage *storage) Put(key, value string) error {
	storage.Lock()
	defer storage.Unlock()

	storage.data[key] = value

	return nil
}

func (storage *storage) Get(key string) (string, error) {
	storage.RLock()
	defer storage.RUnlock()
	value, ok := storage.data[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (storage *storage) Delete(key string) error {
	storage.Lock()
	defer storage.Unlock()

	delete(storage.data, key)
	return nil
}
