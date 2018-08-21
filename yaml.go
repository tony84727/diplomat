package diplomat

import (
	"errors"
	"fmt"
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
	yamlMap = interfaceMapToYAMLMap(root)
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
	var previous YAMLMap
	var current interface{} = yamlMap
	lastID := len(path) - 1
	for i, p := range path {
		switch v := current.(type) {
		case YAMLMap:
			n, exist := v[p]
			if !exist {
				if i != lastID {
					v[p] = make(map[string]interface{})
				}
			}
			previous = v
			current = n
			break
		case string:
			n := map[string]interface{}{p: ""}
			previous[path[i-1]] = n
			previous = n
			current = n[p]
			break
		}
	}
	previous[path[len(path)-1]] = value
	return nil
}
