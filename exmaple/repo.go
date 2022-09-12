package main

import (
	"context"
	model_cache "github.com/mingolm/model-cache"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
)

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
		BackToSource: r.backToSource,
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
