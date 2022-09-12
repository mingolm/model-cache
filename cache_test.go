package model_cache

import (
	"context"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Student struct {
	Id   string
	Name string
}

var ins Cache[string, Student]

func init() {
	ins = New(&Config[string, Student]{
		Marshaler: marshal.JSON,
		Storer:    store.NewMemory(),
		BackToSource: func(ctx context.Context, keys ...string) (results map[string]Student, err error) {
			results = make(map[string]Student, 0)
			for _, key := range keys {
				if row, ok := students[key]; ok {
					results[key] = row
				}
			}
			return results, nil
		},
		KeyPrefix:  "st:",
		Expiration: time.Minute,
	})
}

var students = map[string]Student{
	"1": {Id: "1", Name: "name1"},
	"2": {Id: "2", Name: "name2"},
	"3": {Id: "3", Name: "name3"},
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	row, err := ins.Get(ctx, "1")
	a.Equal(nil, err)
	a.Equal("name1", row.Name)

	// not found
	_, err = ins.Get(ctx, "4")
	a.Equal(store.ErrNotFound, err)
}

func TestMGet(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	rows, err := ins.MGet(ctx, "1", "2", "3", "4", "5")
	a.Equal(nil, err)
	a.Equal(3, len(rows))
}
