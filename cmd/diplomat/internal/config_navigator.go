package internal

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"reflect"
	"strconv"
	"strings"
)

type ConfigNavigator struct {
	previous reflect.Value
	current reflect.Value
	currentType reflect.Type
	valid bool
}

func NewConfigNavigator(config data.Configuration) *ConfigNavigator {
	navigator := &ConfigNavigator{valid:true}
	navigator.setCurrent(reflect.ValueOf(config))
	return navigator
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
	c.previous = c.current
	c.current = value
	if !c.current.IsValid() {
		c.valid = false
		return
	}
	c.currentType = value.Type()
}

func (c *ConfigNavigator) Get(paths ...string) (interface{}, error) {
	i := 0
	for i < len(paths) {
		segment := paths[i]
		switch c.currentType.Kind() {
		case reflect.Ptr, reflect.Interface:
			c.setCurrent(c.current.Elem())
			continue
		case reflect.Slice:
			index,err := strconv.Atoi(segment)
			if err != nil {
				return nil,fmt.Errorf(`%s is a slice. "%s" is not a valid integer`, strings.Join(paths[:i],"."), segment)
			}
			c.setCurrent(c.current.Index(index))
			i++
		case reflect.Struct:
			searcher := FieldSearcher{c.current}
			index, ok := searcher.Search(segment)
			if !ok {
				return nil, fmt.Errorf("%s doesn't exist, possible values %s", segment, structFieldHint(c.current))
			}
			c.setCurrent(c.current.FieldByIndex(index))
			i++
		case reflect.Map:
			c.setCurrent(c.current.MapIndex(reflect.ValueOf(segment)))
			i++
		}
		if !c.valid {
			return nil, fmt.Errorf("%s doesn't exist. current: %v", segment, c.previous.Interface())
		}
	}

	return c.current.Interface(),nil
}
