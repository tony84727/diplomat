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



func (f FieldSearcher) Search(name string) (value reflect.Value, ok bool) {
	for i := 0 ; i < f.value.NumField(); i++ {
		field := f.value.Type().FieldByIndex([]int{i})
		navigateTag, exist := field.Tag.Lookup("navigate")
		if !exist {
			continue
		}
		if name == navigateTag {
			return f.value.Field(i), true
		}
	}
	// search one level down for embedded struct
	for i := 0; i < f.value.NumField(); i++ {
		fieldValue := reflecthelper.Actual(f.value.Field(i))
		fieldType := fieldValue.Type()
		for j := 0; j < fieldType.NumField(); j++ {
			field := fieldType.Field(j)
			value, exist := field.Tag.Lookup("navigate")
			if !exist {
				continue
			}
			if name == value {
				return fieldValue.Field(j), true
			}
		}
	}
	return reflect.Value{}, false
}
