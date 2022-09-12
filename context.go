package model_cache

import "context"

type contextIgnoreStoreMarker struct{}

var ctxIgnoreStoreKey = contextIgnoreStoreMarker{}

// ContextWithSkipStore 返回的 ctx 作为参数调用 Get/MGet 时，回源内容不写入 cache
func ContextWithSkipStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxIgnoreStoreKey, 1)
}
