# Model-Cache

`model-cache`专门为 model 层实现的缓存中间层，主要功能与 CDN 原理类似，程序运行的流程为

1. 调用 repo.Get 获取数据
2. repo 调用 model-cache 尝试从缓存中读取
3. model-cache 判断是否存在缓存，存在则返回，否则进行下一步
4. model-cache 自动调用回源函数，从 repo 中获取数据，并添加至缓存，另外实现了避免缓存穿透，对于 repo 中不存在的数据，会自动存储空值

此处 repo 相当于 CDN 的源站，`model-cache`支持的主要功能点：

1. 数据库自动回源配置
2. 支持设置缓存过期时间
3. 支持缓存空值，避免缓存穿透
4. 支持自定义序列化，默认 JSON 
5. 支持自定义缓存实例，默认 内存，另外支持 redis
6. 支持泛型

```go
package main

import (
	"context"
	model_cache "github.com/mingolm/model-cache"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
)


func main() {
	// 创建 repo 实例
	userRepo := NewRepo()
	row, err := userRepo.Get(context.Background(), 1001)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("user: %+v\n", row)
}


func NewRepo() *Repo {
	r := &Repo{
		users: []user{
			{Id: 1000, Name: "mingo1000"},
			{Id: 1001, Name: "mingo1001"},
			{Id: 1002, Name: "mingo1002"},
			{Id: 1003, Name: "mingo1003"},
		},
	}
	r.mcache = model_cache.New[uint64, user](&model_cache.Config[uint64, user]{
		Marshaler:    marshal.JSON,
		Storer:       store.NewMemory(),
		BackToSource: r.backToSource, // 自动回源函数
		KeyPrefix:    "user:",
		Expiration:   0,
	})

	return r
}

type Repo struct {
	mcache model_cache.Cache[uint64, user]
	users  []user
}

type user struct {
	Id   uint64
	Name string
}

func (r *Repo) Get(ctx context.Context, id uint64) (row *user, err error) {
	row, err = r.mcache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Repo) backToSource(ctx context.Context, keys ...uint64) (results map[uint64]user, err error) {
	results = make(map[uint64]user, 0)
	for _, key := range keys {
		for _, userRow := range r.users {
			if userRow.Id == key {
				results[userRow.Id] = userRow
			}
		}
	}
	return results, nil
}


```