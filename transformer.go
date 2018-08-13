package diplomat

import (
	"strings"

	"github.com/siongui/gojianfan"
)

type TransformHandler func(nkv NestedKeyValue) NestedKeyValue

const (
	SimplifiedToTranditonal = iota
	TranditionalToSimplified
)

type ChineseTransformer struct {
	mode int
	from string
	to   string
}

func optionToChineseTransformerHandler(o YAMLOption) (TransformHandler, error) {
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
	transformer := ChineseTransformer{
		mode: m,
		from: from.(string),
		to:   to.(string),
	}
	return transformer.getTransformerHandler(), nil
}

func (c ChineseTransformer) getTransformerHandler() TransformHandler {
	return func(nkv NestedKeyValue) NestedKeyValue {
		for _, keys := range nkv.GetKeys() {
			if keys[len(keys)-1] == c.from {
				to := append(keys[:len(keys)-1], c.to)
				if !nkv.HasKey(to...) {
					value, _ := nkv.GetKey(keys...)
					nkv.Set(to, c.transform(value.(string)))
				}
			}
		}
		return nkv
	}
}

func (c ChineseTransformer) transform(in string) string {
	if c.mode == SimplifiedToTranditonal {
		return gojianfan.S2T(in)
	} else {
		return gojianfan.T2S(in)
	}
}
