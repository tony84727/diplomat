package internal

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"reflect"
	"strconv"
	"strings"
)

type ConfigNavigator struct {
	current reflect.Value
	currentType reflect.Type
}

func NewConfigNavigator(config data.Configuration) *ConfigNavigator {
	return &ConfigNavigator{current: reflect.ValueOf(config), currentType: reflect.TypeOf(config)}
}

func structFieldHint(structValue reflect.Value) string {
	keys := make([]string, structValue.NumField())
	for i := 0; i < structValue.NumField() ; i++ {
		v := structValue.FieldByIndex([]int{i})
		keys[i] = v.Type().Name()
	}
	return fmt.Sprintf("%v", keys)
}

func (c *ConfigNavigator) setCurrent(value reflect.Value) {
	c.current = value
	c.currentType = value.Type()
}

func (c *ConfigNavigator) Get(paths ...string) (interface{}, error) {
	i := 0
	for i < len(paths) {
		switch c.currentType.Kind() {
		case reflect.Ptr, reflect.Interface:
			c.setCurrent(c.current.Elem())
			continue
		case reflect.Slice:
			index,err := strconv.Atoi(paths[i])
			if err != nil {
				return nil,fmt.Errorf(`%s is a slice. "%s" is not a valid integer`, strings.Join(paths[:i],"."), paths[i])
			}
			c.setCurrent(c.current.Index(index))
			i++
		case reflect.Struct:
			searcher := FieldSearcher{c.currentType}
			index, ok := searcher.Search(paths[i])
			if !ok {
				return nil, fmt.Errorf("%s doesn't exist, possible values %s", paths[i], structFieldHint(c.current))
			}
			c.setCurrent(c.current.FieldByIndex(index))
			i++
		case reflect.Map:
			c.setCurrent(c.current.MapIndex(reflect.ValueOf(paths[i])))
			i++
		default:
			if i < len(paths)  {
				return nil, fmt.Errorf("%s doesn't exist. current: %T", paths[i], c.current.Interface())
			}
		}
	}

	return c.current.Interface(),nil
}
