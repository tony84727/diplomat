package diplomat

import (
	"strings"

	"github.com/siongui/gojianfan"
)

type PreprocesserFunc func(nkv *NestedKeyValue) error

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

func (c ChineseTransformer) getPreprocessorFunc() PreprocesserFunc {
	return func(nkv *NestedKeyValue) error {
		for _, keys := range nkv.GetKeys() {
			if keys[len(keys)-1] == c.from {
				to := make([]string, len(keys))
				lastID := len(to) - 1
				for i := 0; i < lastID; i++ {
					to[i] = keys[i]
				}
				to[lastID] = c.to
				value, exist := nkv.GetKey(keys...)
				if exist {
					nkv.Set(to, c.transform(value.(string)))
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
	return func(nkv *NestedKeyValue) error {
		for _, v := range preprocessors {
			err := v(nkv)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type PreprocesserFuncFactory func(YAMLOption) (PreprocesserFunc, error)

type preprocessorManager map[string]PreprocesserFuncFactory

var preprocessorManagerInstance preprocessorManager = make(map[string]PreprocesserFuncFactory)

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
