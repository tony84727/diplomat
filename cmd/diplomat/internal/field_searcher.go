package internal

import (
	"github.com/tony84727/diplomat/pkg/reflecthelper"
	"reflect"
)

// FieldSearcher is used for search field of a reflect type
// with "navigate" tag
type FieldSearcher struct {
	value reflect.Value
}



func (f FieldSearcher) Search(name string) (fieldIndex []int, ok bool) {
	for i := 0 ; i < f.value.NumField(); i++ {
		f := f.value.Type().FieldByIndex([]int{i})
		value, exist := f.Tag.Lookup("navigate")
		if !exist {
			continue
		}
		if name == value {
			return []int{i}, true
		}
	}
	// search one level down for embedded struct
	for i := 0; i < f.value.NumField(); i++ {
		fieldValue := reflecthelper.Actual(f.value.Field(i))
		top := fieldValue.Type()
		for j := 0; j < top.NumField(); j++ {
			f := top.FieldByIndex([]int{j})
			value, exist := f.Tag.Lookup("navigate")
			if !exist {
				continue
			}
			if name == value {
				return []int{i,j}, true
			}
		}
	}
	return nil, false
}
