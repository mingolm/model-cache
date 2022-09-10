package marshal

import "reflect"

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	String() string
}

// UnmarshalIntoArray 判断 pointer of array
// bytes 是 nil
// 如果是基础类型
func UnmarshalIntoArray(bytesArray [][]byte, values interface{}, unmarshalFunc func(data []byte, v interface{}) error) (err error) {
	indirectValue := reflect.Indirect(reflect.ValueOf(values))

	sliceEle := indirectValue
	l := len(bytesArray)
	if sliceEle.CanSet() && sliceEle.Cap() < l {
		sliceEle = reflect.MakeSlice(sliceEle.Type(), l, l)
		indirectValue.Set(sliceEle)
	} else if sliceEle.Len() < l {
		sliceEle.SetLen(l)
	}
	for i, bs := range bytesArray {
		if bs == nil {
			continue
		}
		e := sliceEle.Index(i)
		if e.Kind() != reflect.Ptr {
			e = e.Addr()
		} else if e.IsNil() {
			e.Set(reflect.New(e.Type().Elem()))
		}

		v := e.Interface()
		err = unmarshalFunc(bs, v)
		if err != nil {
			return err
		}
	}

	return nil
}
