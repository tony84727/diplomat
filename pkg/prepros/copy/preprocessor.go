package copy

import (
	"errors"
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/prepros"
	"github.com/tony84727/diplomat/pkg/prepros/types"
)

type Preprocessor struct {
}

func (p Preprocessor) Process(translation data.Translation, option interface{}) error {
	config, err := p.parseConfig(option)
	if err != nil {
		return err
	}
	if err := p.validConfig(config); err != nil {
		return err
	}
	walker := data.NewTranslationWalker(translation)
	return walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
		if paths[len(paths)-1] == config.From {
			if c := textNode.GetParent().GetChild(config.To); c == nil {
				toNode := data.NewTranslation(config.To)
				toNode.SetText(*textNode.GetText())
				textNode.GetParent().AddChild(toNode)
			}
		}
		return nil
	})
}

type Config struct {
	From string
	To   string
}

func (Preprocessor) parseConfig(option interface{}) (*Config, error) {
	m, ok := option.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("expect option to be a map, go %v", option)
	}
	from, ok := m["from"].(string)
	if !ok {
		return nil, errors.New("expecting from option")
	}
	to, ok := m["to"].(string)
	if !ok {
		return nil, errors.New("expecting to option")
	}
	return &Config{from, to}, nil
}

func (Preprocessor) validConfig(config *Config) error {
	if len(config.From) <= 0 {
		return errors.New("from is empty")
	}
	if len(config.To) <= 0 {
		return errors.New("to is empty")
	}
	return nil
}

func New() types.Preprocessor {
	return &Preprocessor{}
}

func init() {
	prepros.GlobalRegistry.Registry("copy", types.FactoryWrapper{Constructor: New})
}
