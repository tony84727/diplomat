package diplomat

import (
	"errors"
	"fmt"

	set "github.com/deckarep/golang-set"
)

func interfaceMapToStringMap(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, i := range in {
		switch v := i.(type) {
		case map[interface{}]interface{}:
			out[k.(string)] = interfaceMapToStringMap(v)
			break
		case []interface{}:
			out[k.(string)] = checkSliceForStringMap(v)
			break
		default:
			out[k.(string)] = i
		}
	}
	return out
}

func checkSliceForStringMap(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		switch v := e.(type) {
		case map[interface{}]interface{}:
			out[i] = interfaceMapToStringMap(v)
			break
		case []interface{}:
			out[i] = checkSliceForStringMap(v)
			break
		default:
			out[i] = e
		}
	}
	return out
}

type YAMLOption struct {
	data interface{}
}

func (y YAMLOption) IsSlice(paths ...interface{}) (bool, error) {
	child, err := y.Get(paths...)
	if err != nil {
		return false, err
	}
	_, ok := child.([]interface{})
	return ok, nil
}

func (y YAMLOption) Len(paths ...interface{}) (int, error) {
	child, err := y.Get(paths...)
	if err != nil {
		return 0, err
	}
	switch v := child.(type) {
	case map[string]interface{}:
	case []interface{}:
		return len(v), nil
	}
	return 0, errors.New("the element is not a map nor a slice")
}

func (y YAMLOption) Get(paths ...interface{}) (interface{}, error) {
	current := y.data
	for i, p := range paths {
		switch v := p.(type) {
		case string:
			cv, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%v is not a map", paths[:i])
			}
			current = cv[v]
			break
		case int:
			cv, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%v is not a slice", paths[:i])
			}
			current = cv[v]
		}
	}
	return current, nil
}

func (y *YAMLOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var varient interface{}
	err := unmarshal(&varient)
	if err != nil {
		return err
	}
	switch v := varient.(type) {
	case map[interface{}]interface{}:
		y.data = interfaceMapToStringMap(v)
		break

	case []interface{}:
		y.data = checkSliceForStringMap(v)
		break
	default:
		return fmt.Errorf("YAMLOption should be either map or slice")
	}
	return nil
}
func interfaceMapToYAMLMap(in map[interface{}]interface{}) YAMLMap {
	out := make(YAMLMap, len(in))
	for k, i := range in {
		switch v := i.(type) {
		case map[interface{}]interface{}:
			out[k.(string)] = interfaceMapToYAMLMap(v)
			break
		case []interface{}:
			out[k.(string)] = checkSliceForYAMLMap(v)
			break
		default:
			out[k.(string)] = i
		}
	}
	return out
}

func checkSliceForYAMLMap(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		switch v := e.(type) {
		case map[interface{}]interface{}:
			out[i] = interfaceMapToYAMLMap(v)
			break
		case []interface{}:
			out[i] = checkSliceForYAMLMap(v)
			break
		default:
			out[i] = e
		}
	}
	return out
}

// YAMLMap is a tree node to store nested translations.
type YAMLMap map[string]interface{}

func (yamlMap YAMLMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var root map[interface{}]interface{}
	err := unmarshal(&root)
	if err != nil {
		return err
	}
	data := interfaceMapToYAMLMap(root)
	for k, v := range data {
		yamlMap[k] = v
	}
	return nil
}

func (yamlMap YAMLMap) GetKey(paths ...string) (value interface{}, exist bool) {
	if len(paths) <= 0 {
		return yamlMap, true
	}
	d, exist := yamlMap[paths[0]]
	if !exist {
		return nil, false
	}
	v, ok := d.(string)
	if ok {
		return v, true
	}
	return d.(YAMLMap).GetKey(paths[1:]...)
}

func (yamlMap YAMLMap) GetKeys() [][]string {
	keys := make([][]string, 0, 1)
	for k, v := range yamlMap {
		switch i := v.(type) {
		case string:
			keys = append(keys, []string{k})
			break
		case YAMLMap:
			nKeys := i.GetKeys()
			for _, s := range nKeys {
				keys = append(keys, append([]string{k}, s...))
			}
			break
		default:
			continue
		}
	}
	return keys
}

func (yamlMap YAMLMap) filterBySelectorOnBase(base []string, s Selector) YAMLMap {
	filtered := make(map[string]interface{})
	for k, i := range yamlMap {
		key := make([]string, len(base)+1)
		for b, bk := range base {
			key[b] = bk
		}
		key[len(base)] = k
		switch v := i.(type) {
		case string:
			if s.IsValid(key) {
				filtered[k] = v
			}
			break
		case YAMLMap:
			if s.IsValid(key) {
				filtered[k] = v
			} else {
				child := v.filterBySelectorOnBase(key, s)
				if len(child) > 0 {
					filtered[k] = child
				}
			}
		}
	}
	return filtered
}

func (yamlMap YAMLMap) FilterBySelector(s Selector) YAMLMap {
	return yamlMap.filterBySelectorOnBase([]string{}, s)
}

func (yamlMap YAMLMap) HasKey(keys ...string) bool {
	_, exist := yamlMap.GetKey(keys...)
	return exist
}

func (yamlMap YAMLMap) LanguageHasKey(language string, keys ...string) bool {
	keys = append(keys, language)
	return yamlMap.HasKey(keys...)
}

func (yamlMap YAMLMap) Set(path []string, value string) error {
	if len(path) < 1 {
		return errors.New("except at least on path")
	}
	current := yamlMap
	for i := 0; i < len(path); i++ {
		p := path[i]
		next, exist := current[p]
		if exist {
			if m, isMap := next.(YAMLMap); isMap {
				current = m
				continue
			}
		}
		current[p] = make(YAMLMap)
		if i == len(path)-1 {
			current[p] = value
			continue
		}
		current = current[p].(YAMLMap)
	}
	return nil
}

func (yamlMap YAMLMap) GetLanguages() []string {
	s := set.NewSet()
	for _, k := range yamlMap.GetKeys() {
		s.Add(k[len(k)-1])
	}
	list := s.ToSlice()
	stringList := make([]string, len(list))
	for i, language := range list {
		stringList[i] = language.(string)
	}
	return stringList
}

func MergeYAMLMaps(maps ...YAMLMap) YAMLMap {
	all := make(YAMLMap)
	for _, s := range maps {
		for _, k := range s.GetKeys() {
			v, _ := s.GetKey(k...)
			all.Set(k, v.(string))
		}
	}
	return all
}
