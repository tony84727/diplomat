package transfrom

import (
	"github.com/insufficientchocolate/diplomat/pkg/data"
	"github.com/siongui/gojianfan"
)

type ChineseTransformMode int

const (
	SimplifiedToTranditonal ChineseTransformMode = iota
	TranditionalToSimplified
)

type ChineseTransformerOption struct {
	mode ChineseTransformMode
	from string
	to   string
}

type ChineseTransformer struct {
	option ChineseTransformerOption
}

func (c ChineseTransformer) Transform(translation data.Translation) error {
	walker := data.NewTranslationWalker(translation)
	toAdd := make([]data.Translation, 0)
	walker.ForEachTextNode(func(path string, textNode data.Translation) error {
		if textNode.GetKey() == c.option.from {
			toAdd = append(toAdd, textNode)
		}
		return nil
	})
	for _, node := range toAdd {
		translated := c.translate(*node.GetText())
		translatedNode :=data.NewTranslation(c.option.to)
		translatedNode.SetText(translated)
		node.GetParent().AddChild(translatedNode)
	}
	return nil
}

func (c ChineseTransformer) translate(input string) string {
	if c.option.mode == SimplifiedToTranditonal {
		return gojianfan.S2T(input)
	}
	return gojianfan.T2S(input)
}

func NewChineseTransformer(option ChineseTransformerOption) *ChineseTransformer {
	return &ChineseTransformer{option}
}
