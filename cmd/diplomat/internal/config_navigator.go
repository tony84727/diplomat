package internal

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"reflect"
	"strconv"
	"strings"
)

type ConfigNavigator struct {
	config data.Configuration
}

func NewConfigNavigator(config data.Configuration) *ConfigNavigator {
	return &ConfigNavigator{config: config}
}

func structFieldHint(structValue reflect.Value) string {
	keys := make([]string, structValue.NumField())
	for i := 0; i < structValue.NumField() ; i++ {
		v := structValue.FieldByIndex([]int{i})
		keys[i] = v.Type().Name()
	}
	return fmt.Sprintf("%v", keys)
}

// searchField given the field name, search the name by "navigate" tag.
func searchField(fieldName string, theType reflect.Type) (index []int, ok bool) {
	for i := 0 ; i < theType.NumField(); i++ {
		f := theType.FieldByIndex([]int{i})
		value, exist := f.Tag.Lookup("navigate")
		if !exist {
			continue
		}
		if fieldName == value {
			return []int{i}, true
		}
	}
	// search one level down for embedded struct
	for i := 0; i < theType.NumField(); i++ {
		top := theType.FieldByIndex([]int{i})
		for j := 0; j < top.Type.NumField(); j++ {
			f := top.Type.FieldByIndex([]int{j})
			value, exist := f.Tag.Lookup("navigate")
			if !exist {
				continue
			}
			if fieldName == value {
				return []int{i,j}, true
			}
		}
	}
	return nil, false
}

func (c ConfigNavigator) Get(paths ...string) (interface{}, error) {
	currentValue := reflect.ValueOf(c.config)
	i := 0
	for i < len(paths) {
		currentType := currentValue.Type()
		switch currentType.Kind() {
		case reflect.Ptr, reflect.Interface:
			currentValue = currentValue.Elem()
		case reflect.Slice:
			index,err := strconv.Atoi(paths[i])
			if err != nil {
				return nil,fmt.Errorf(`%s is a slice. "%s" is not a valid integer`, strings.Join(paths[:i],"."), paths[i])
			}
			currentValue = currentValue.Index(index)
			i++
		case reflect.Struct:
			index, ok := searchField(paths[i], currentType)
			if !ok {
				return nil, fmt.Errorf("%s doesn't exist, possible values %s", paths[i], structFieldHint(currentValue))
			}
			currentValue = currentValue.FieldByIndex(index)
			i++
		case reflect.Map:
			currentValue = currentValue.MapIndex(reflect.ValueOf(paths[i]))
			i++
		default:
			if i < len(paths)  {
				return nil, fmt.Errorf("%s doesn't exist. current: %T", paths[i], currentValue.Interface())
			}
		}
	}

	return currentValue.Interface(),nil
}
