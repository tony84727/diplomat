package diplomat

import (
	"errors"
	"strings"

	"github.com/siongui/gojianfan"
)

type PreprocesserFunc func(yamlMap YAMLMap) error

const (
	SimplifiedToTranditonal = iota
	TranditionalToSimplified
)

type ChineseTransformer struct {
	mode int
	from string
	to   string
}

func optionToChineseTransformer(o YAMLOption) (*ChineseTransformer, error) {
	mode, err := o.Get("mode")
	if err != nil {
		return nil, err
	}
	from, err := o.Get("from")
	if err != nil {
		return nil, err
	}
	to, err := o.Get("to")
	if err != nil {
		return nil, err
	}
	m := SimplifiedToTranditonal
	if strings.ToLower(mode.(string)) != "s2t" {
		m = TranditionalToSimplified
	}
	return &ChineseTransformer{
		mode: m,
		from: from.(string),
		to:   to.(string),
	}, nil
}

func optionToChinesePreprocessorFunc(o YAMLOption) (PreprocesserFunc, error) {
	transformer, err := optionToChineseTransformer(o)
	if err != nil {
		return nil, err
	}
	return transformer.getPreprocessorFunc(), nil
}

func chinesePreprocessorFactory(o YAMLOption) (PreprocesserFunc, error) {
	isSlice, err := o.IsSlice()
	if err != nil {
		return nil, err
	}
	if !isSlice {
		return nil, errors.New("expect array")
	}
	l, err := o.Len()
	if err != nil {
		return nil, err
	}
	funcs := make([]PreprocesserFunc, l)
	options, err := o.Get()
	if err != nil {
		return nil, err
	}
	for i, op := range options.([]interface{}) {
		f, err := optionToChinesePreprocessorFunc(YAMLOption{data: op})
		if err != nil {
			return nil, err
		}
		funcs[i] = f
	}
	return combinePreprocessor(funcs...), nil
}

func (c ChineseTransformer) getPreprocessorFunc() PreprocesserFunc {
	return func(yamlMap YAMLMap) error {
		for _, keys := range yamlMap.GetKeys() {
			if keys[len(keys)-1] == c.from {
				to := make([]string, len(keys))
				lastID := len(to) - 1
				for i := 0; i < lastID; i++ {
					to[i] = keys[i]
				}
				to[lastID] = c.to
				value, exist := yamlMap.GetKey(keys...)
				if exist {
					yamlMap.Set(to, c.transform(value.(string)))
				}
			}
		}
		return nil
	}
}

func (c ChineseTransformer) transform(in string) string {
	if c.mode == SimplifiedToTranditonal {
		return gojianfan.S2T(in)
	}
	return gojianfan.T2S(in)
}

func combinePreprocessor(preprocessors ...PreprocesserFunc) PreprocesserFunc {
	return func(yamlMap YAMLMap) error {
		for _, v := range preprocessors {
			err := v(yamlMap)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type PreprocesserFuncFactory func(YAMLOption) (PreprocesserFunc, error)

type preprocessorManager map[string]PreprocesserFuncFactory

var preprocessorManagerInstance preprocessorManager

func init() {
	preprocessorManagerInstance = make(map[string]PreprocesserFuncFactory)
	RegisterPreprocessorFuncFactory("chinese", chinesePreprocessorFactory)
}

func (pm preprocessorManager) buildPreprocessors(configs []PreprocessorConfig) ([]PreprocesserFunc, error) {
	preprocessors := make([]PreprocesserFunc, 0, len(configs))
	for _, c := range configs {
		factory, exists := pm[c.Type]
		if exists {
			f, err := factory(c.Options)
			if err != nil {
				return nil, err
			}
			preprocessors = append(preprocessors, f)
		}
	}
	return preprocessors, nil
}

func RegisterPreprocessorFuncFactory(preprocessorType string, factory PreprocesserFuncFactory) {
	preprocessorManagerInstance[preprocessorType] = factory
}
