package model_cache

import (
	"context"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
	"time"
)

type Cache[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Get(ctx context.Context, key string) (T, error)
	MGet(ctx context.Context, keys ...string) (*[]T, error)
	Del(ctx context.Context, keys ...string) error
}

func New[T any](marshaler marshal.Marshaler, storer store.Storer) Cache[T] {
	return &cache[T]{
		marshaler: marshaler,
		storer:    storer,
	}
}

type cache[T any] struct {
	marshaler marshal.Marshaler
	storer    store.Storer
}

func (c *cache[T]) Set(ctx context.Context, key string, value T, expiration time.Duration) (err error) {
	bs, err := c.marshaler.Marshal(value)
	if err != nil {
		return err
	}
	return c.storer.Set(ctx, key, bs, expiration)
}

func (c *cache[T]) Get(ctx context.Context, key string) (row T, err error) {
	bs, err := c.storer.Get(ctx, key)
	if err != nil {
		return row, err
	}
	if err = c.marshaler.Unmarshal(bs, &row); err != nil {
		return row, err
	}
	return row, nil
}

func (c *cache[T]) MGet(ctx context.Context, keys ...string) (rows *[]T, err error) {
	bsArray, err := c.storer.MGet(ctx, keys...)
	if err != nil {
		return nil, err
	}
	if err = marshal.UnmarshalIntoArray(bsArray, &rows, c.marshaler.Unmarshal); err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *cache[T]) Del(ctx context.Context, keys ...string) (err error) {
	if err := c.storer.Del(ctx, keys...); err != nil {
		return err
	}
	return nil
}
