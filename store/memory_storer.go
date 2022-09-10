package store

import (
	"context"
	"sync"
	"time"
)

func NewMemory() Storer {
	return &Memory{
		db:        sync.Map{},
		dbExpired: sync.Map{},
	}
}

type Memory struct {
	db        sync.Map
	dbExpired sync.Map
}

func (m *Memory) Get(ctx context.Context, key string) (bs []byte, err error) {
	v, ok := m.db.Load(key)
	if !ok {
		return nil, ErrNotFound
	}

	if _, expired := m.filterExpired(key); expired {
		return nil, ErrNotFound
	}

	return v.([]byte), nil
}

func (m *Memory) MGet(ctx context.Context, keys ...string) (results [][]byte, err error) {
	for _, key := range keys {
		v, ok := m.db.Load(key)
		if ok {
			if _, expired := m.filterExpired(key); expired {
				v = nil
			}
		}
		if v == nil {
			results = append(results, nil)
		} else {
			results = append(results, v.([]byte))
		}
	}

	return results, nil
}

func (m *Memory) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	m.db.Store(key, value)
	if expiration != 0 {
		m.dbExpired.Store(key, time.Now().Add(expiration))
	}
	return nil
}

func (m *Memory) Del(ctx context.Context, keys ...string) (err error) {
	for _, key := range keys {
		m.db.Delete(key)
		m.dbExpired.Delete(key)
	}

	return nil
}

func (m *Memory) String() string {
	return "memory"
}

func (m *Memory) filterExpired(key string) (ttl time.Duration, expired bool) {
	e, ok := m.dbExpired.Load(key)
	if !ok {
		return -1, false
	}

	expiredTime := e.(time.Time)

	if time.Now().After(expiredTime) {
		m.db.Delete(key)
		m.dbExpired.Delete(key)
		return 0, true
	}

	return expiredTime.Sub(time.Now()), false
}
