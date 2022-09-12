package model_cache

import (
	"bytes"
)

var emptyCacheBsValue = []byte("{}")

func isEmptyValue(v []byte) bool {
	return bytes.Equal(emptyCacheBsValue, v)
}
