package yaml

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"gopkg.in/yaml.v2"
)

type TranslationParser struct {
	content []byte
	root    data.Translation
}

type translationFile yaml.MapSlice

func (p *TranslationParser) GetTranslation() (data.Translation, error) {
	if p.root != nil {
		return p.root, nil
	}
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return p.root, nil
}

func (p *TranslationParser) parse() error {
	var translations translationFile
	err := yaml.Unmarshal(p.content, &translations)
	if err != nil {
		return err
	}

	root := data.NewTranslation("")
	err = p.assignTranslations(root, translations)
	if err != nil {
		return err
	}
	p.root = root
	return nil
}

func (p TranslationParser) assignTranslations(root data.Translation, input translationFile) error {
	for _, item := range input {
		stringKey, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("unexpected %v", input)
		}
		current := data.NewTranslation(stringKey)
		switch v := item.Value.(type) {
		case translationFile:
			p.assignTranslations(current, translationFile(v))
		case string:
			current.SetText(v)
		default:
			return fmt.Errorf("unexpected %v(%T)", v, v)
		}
		root.AddChild(current)
	}
	return nil
}

func NewParser(content []byte) *TranslationParser {
	return &TranslationParser{content: content}
}
