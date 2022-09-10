package store

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type StorerType string

type Storer interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	MGet(ctx context.Context, keys ...string) ([][]byte, error)
	Del(ctx context.Context, keys ...string) error
	String() string
}

var (
	ErrNotFound = errors.New("cache not found")
)

func errSwap(err error) error {
	switch err {
	case redis.Nil:
		return ErrNotFound
	default:
		return err
	}
}
