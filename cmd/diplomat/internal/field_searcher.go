package internal

import "reflect"

// FieldSearcher is used for search field of a reflect type
// with "navigate" tag
type FieldSearcher struct {
	theType reflect.Type
}

func (f FieldSearcher) Search(name string) (fieldIndex []int, ok bool) {
	for i := 0 ; i < f.theType.NumField(); i++ {
		f := f.theType.FieldByIndex([]int{i})
		value, exist := f.Tag.Lookup("navigate")
		if !exist {
			continue
		}
		if name == value {
			return []int{i}, true
		}
	}
	// search one level down for embedded struct
	for i := 0; i < f.theType.NumField(); i++ {
		top := f.theType.FieldByIndex([]int{i})
		for j := 0; j < top.Type.NumField(); j++ {
			f := top.Type.FieldByIndex([]int{j})
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
