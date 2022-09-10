package model_cache

import (
	"context"
	"github.com/mingolm/model-cache/marshal"
	"github.com/mingolm/model-cache/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Student struct {
	Id   uint64
	Name string
	Age  uint8
}

var ins Cache[Student]

func init() {
	ins = New[Student](marshal.JSON, store.NewMemory())
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	err := ins.Set(ctx, "test", Student{
		Id:   1,
		Name: "test1",
		Age:  18,
	}, 0)
	a.Equal(nil, err)
	row, err := ins.Get(ctx, "test")
	a.Equal(nil, err)
	a.Equal("test1", row.Name)
}

func TestMGet(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	for _, row := range []Student{
		{
			Id:   1,
			Name: "test1",
			Age:  18,
		},
		{
			Id:   2,
			Name: "test2",
			Age:  19,
		},
		{
			Id:   3,
			Name: "test3",
			Age:  20,
		},
	} {
		err := ins.Set(ctx, row.Name, row, 0)
		a.Equal(nil, err)
	}

	row, err := ins.Get(ctx, "test1")
	a.Equal(nil, err)
	a.Equal("test1", row.Name)

	row, err = ins.Get(ctx, "test2")
	a.Equal(nil, err)
	a.Equal("test2", row.Name)
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	err := ins.Set(ctx, "test", Student{
		Id:   1,
		Name: "test1",
		Age:  18,
	}, 0)
	a.Equal(nil, err)
	row, err := ins.Get(ctx, "test")
	a.Equal(nil, err)
	a.Equal("test1", row.Name)

	err = ins.Del(ctx, "test")
	a.Equal(nil, err)
	row, err = ins.Get(ctx, "test")
	a.Equal(store.ErrNotFound, err)
}
