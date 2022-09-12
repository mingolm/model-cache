package model_cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
	"github.com/spf13/cast"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (*V, error)
	MGet(ctx context.Context, keys ...K) ([]V, error)
	Del(ctx context.Context, keys ...K) error
	Refresh(ctx context.Context, keys ...K) (map[K]V, error)
}

func New[K comparable, V any](conf *Config[K, V]) Cache[K, V] {
	return &cache[K, V]{
		conf,
	}
}

type cache[K comparable, V any] struct {
	*Config[K, V]
}

type Config[K comparable, V any] struct {
	Marshaler marshal.Marshaler
	Storer    store.Storer
	// 回源函数
	// 返回值 results 为 key => row，map key 不存在表示对应的 row 为 notfound
	BackToSource BackToSource[K, V]
	// 禁止缓存空值（BackToSource 未返回值或者返回 nil）
	DisableEmptyValueCache bool
	// 缓存前缀
	KeyPrefix string
	// 缓存过期时间
	Expiration time.Duration
}

type BackToSource[K comparable, V any] func(ctx context.Context, keys ...K) (results map[K]V, err error)

func (c *cache[K, V]) Get(ctx context.Context, key K) (row *V, err error) {
	keyStr := c.getCacheKey(key)
	bs, err := c.Storer.Get(ctx, keyStr)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return nil, err
		}
		results, err := c.backSourceAndSetStore(ctx, key)
		if err != nil {
			return nil, err
		}
		tRow, ok := results[key]
		if ok {
			return &tRow, nil
		}

		return nil, store.ErrNotFound
	}
	if isEmptyValue(bs) {
		return nil, store.ErrNotFound
	}
	if err = c.Marshaler.Unmarshal(bs, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (c *cache[K, V]) MGet(ctx context.Context, keys ...K) (rows []V, err error) {
	bsArray, err := c.Storer.MGet(ctx, c.getCacheKeys(keys...)...)
	if err != nil {
		return nil, err
	}
	var missCacheKeys []K
	var statusCacheKeyIndexMap = make(map[K]int8, len(keys))
	for ib, bs := range bsArray {
		if bs == nil {
			missCacheKeys = append(missCacheKeys, keys[ib])
			statusCacheKeyIndexMap[keys[ib]] = -1
		} else if isEmptyValue(bs) {
			statusCacheKeyIndexMap[keys[ib]] = 0
		} else {
			statusCacheKeyIndexMap[keys[ib]] = 1
		}
	}
	if err = marshal.UnmarshalIntoArray(bsArray, &rows, c.Marshaler.Unmarshal); err != nil {
		return nil, err
	}

	missResults, err := c.backSourceAndSetStore(ctx, missCacheKeys...)
	if err != nil {
		return nil, err
	}

	// 保证返回顺序
	var (
		results  []V
		rowIndex int
	)
	for _, key := range keys {
		switch statusCacheKeyIndexMap[key] {
		case -1:
			if missRow, ok := missResults[key]; ok {
				results = append(results, missRow)
			}
		case 0:
		case 1:
			results = append(results, rows[rowIndex])
			rowIndex++
		}
	}

	return results, nil
}

func (c *cache[K, V]) Del(ctx context.Context, keys ...K) (err error) {
	if err := c.Storer.Del(ctx, c.getCacheKeys(keys...)...); err != nil {
		return err
	}
	return nil
}

func (c *cache[K, V]) Refresh(ctx context.Context, keys ...K) (results map[K]V, err error) {
	results, err = c.backSourceAndSetStore(ctx, keys...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *cache[K, V]) getCacheKeys(keys ...K) (cacheKeys []string) {
	for _, key := range keys {
		cacheKeys = append(cacheKeys, c.getCacheKey(key))
	}

	return cacheKeys
}

func (c *cache[K, V]) getCacheKey(key K) (cacheKey string) {
	keyStr, err := cast.ToStringE(key)
	if err != nil {
		panic(fmt.Errorf("key to string failed: %w", err))
	}

	return keyStr
}

func (c *cache[K, V]) backSourceAndSetStore(ctx context.Context, keys ...K) (results map[K]V, err error) {
	results, err = c.BackToSource(ctx, keys...)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		var bs []byte
		if val, ok := results[key]; ok {
			bs, err = c.Marshaler.Marshal(val)
			if err != nil {
				return nil, err
			}
		}
		// 禁止缓存空值
		if bs == nil && !c.disableStoreEmptyValue(ctx) {
			bs = emptyCacheBsValue
		}
		if bs != nil {
			if err = c.Storer.Set(ctx, c.getCacheKey(key), bs, c.Expiration); err != nil {
				return nil, err
			}
		}
	}

	return results, nil
}

func (c *cache[K, V]) disableStoreEmptyValue(ctx context.Context) bool {
	return c.DisableEmptyValueCache || ctx.Value(ctxIgnoreStoreKey) != nil
}
