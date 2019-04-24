package chinese

import (
	"errors"
	"fmt"
	"github.com/siongui/gojianfan"
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/prepros"
)

type TransformMode int

const (
	SimplifiedToTranditonal TransformMode = iota
	TranditionalToSimplified
)

type Config struct {
	Mode TransformMode
	From string
	To   string
}

type Preprocessor struct {
}

func (Preprocessor) parseConfig(option interface{}) (*Config, error) {
	m, ok := option.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("expect option to be a map, got %v", option)
	}
	mode, ok := m["mode"].(string)
	if !ok {
		return nil, errors.New("expecting mode option: s2t or t2s")
	}
	from, ok := m["from"].(string)
	if !ok {
		return nil, errors.New("expecting from option")
	}
	to, ok := m["to"].(string)
	if !ok {
		return nil, errors.New("expecting to option")
	}
	enumMode := SimplifiedToTranditonal
	if mode == "t2s" {
		enumMode = TranditionalToSimplified
	}
	return &Config{
		Mode: enumMode,
		From: from,
		To:   to,
	}, nil
}

func (Preprocessor) transform(config *Config, input string) string {
	if config.Mode == SimplifiedToTranditonal {
		return gojianfan.S2T(input)
	}
	return gojianfan.T2S(input)
}

func (p Preprocessor) Process(translation data.Translation, option interface{}) error {
	config, err := p.parseConfig(option)
	if err != nil {
		return err
	}
	walker := data.NewTranslationWalker(translation)
	return walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
		if paths[len(paths)-1] == config.From {
			keyNode := textNode.GetParent()
			// ignore if "to" node already exists
			if c := keyNode.GetChild(config.To); c != nil {
				return nil
			}
			toNode := data.NewTranslation(config.To)
			toNode.SetText(p.transform(config, *textNode.GetText()))
			keyNode.AddChild(toNode)
		}
		return nil
	})
}

func init() {
	prepros.GlobalRegistry.Registry("chinese", &Preprocessor{})
}
