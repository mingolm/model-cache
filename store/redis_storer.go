package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedis(cli redis.Cmdable) Storer {
	return &Redis{
		client: cli,
	}
}

type Redis struct {
	client redis.Cmdable
}

func (st *Redis) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	err := st.client.Set(ctx, key, value, expiration).Err()
	return errSwap(err)
}

func (st *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	bs, err := st.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errSwap(err)
	}
	return bs, nil
}

func (st *Redis) MGet(ctx context.Context, keys ...string) ([][]byte, error) {
	result, err := st.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, errSwap(err)
	}

	values := make([][]byte, len(keys), len(keys))
	for idx, val := range result {
		if val == nil {
			values[idx] = nil
		} else {
			values[idx] = []byte(val.(string))
		}
	}

	return values, err
}

func (st *Redis) Del(ctx context.Context, keys ...string) error {
	err := st.client.Del(ctx, keys...).Err()
	return errSwap(err)
}

func (st *Redis) String() string {
	return "redis"
}
