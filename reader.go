package diplomat

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/go-yaml/yaml"
)

type YAMLOption struct {
	data interface{}
}

func (y YAMLOption) Get(paths ...interface{}) (interface{}, error) {
	current := y.data
	for i, p := range paths {
		switch v := p.(type) {
		case string:
			cv, ok := current.(map[interface{}]interface{})
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
		return nil
	}
	y.data = varient
	return nil
}

type PreprocessorConfig struct {
	Type    string
	Options YAMLOption
}

type OutputConfig struct {
	Selectors []string
	Template  MessengerConfig
}

// Outline is the struct of translation file.
type Outline struct {
	Version       string
	Preprocessors []PreprocessorConfig
	Output        []OutputConfig
}

func Read(path string) (*Outline, error) {
	var content Outline
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// NestedKeyValue is a tree node to store nested translations.
type NestedKeyValue struct {
	data map[string]interface{}
}

func nkvDataFromStringMap(m map[string]interface{}) (map[string]interface{}, error) {
	for k, p := range m {
		switch v := p.(type) {
		case int:
			m[k] = strconv.Itoa(v)
			break
		case string:
			m[k] = v
		case map[interface{}]interface{}:
			stringMap := make(map[string]interface{}, len(v))
			for i, j := range v {
				stringMap[i.(string)] = j
			}
			anotherNkv, err := nkvDataFromStringMap(stringMap)
			if err != nil {
				return m, err
			}
			m[k] = NestedKeyValue{anotherNkv}
			break
		default:
			return m, fmt.Errorf("unexcepted type: %T at %s", v, k)
		}
	}
	return m, nil
}

func (nkv *NestedKeyValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var root map[string]interface{}
	err := unmarshal(&root)
	if err != nil {
		return err
	}
	d, err := nkvDataFromStringMap(root)
	if err != nil {
		return err
	}
	nkv.data = d
	return nil
}

func (nkv NestedKeyValue) GetKey(paths ...string) (value interface{}, exist bool) {
	if len(paths) <= 0 {
		return nkv, true
	}
	d, exist := nkv.data[paths[0]]
	if !exist {
		return nil, false
	}
	v, ok := d.(string)
	if ok {
		return v, true
	}
	return d.(NestedKeyValue).GetKey(paths[1:]...)
}

func (nkv NestedKeyValue) GetKeys() [][]string {
	keys := make([][]string, 0, 1)
	for k, v := range nkv.data {
		switch i := v.(type) {
		case string:
			keys = append(keys, []string{k})
			break
		case NestedKeyValue:
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

func (nkv *NestedKeyValue) Put(value string, paths ...string) {
	if len(paths) <= 0 {
		panic("paths is empty")
	}

}

type PartialTranslation struct {
	data map[string]NestedKeyValue
}
