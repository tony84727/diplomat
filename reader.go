package diplomat

import (
	"fmt"
	"io/ioutil"

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
