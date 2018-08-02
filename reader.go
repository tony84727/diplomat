package diplomat

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Translation presents languageCode => translatedText mapping
type Translation = map[string]string

// Translations presents translationKey => translation (for different language) mapping
type Translations = map[string]Translation

// Fragment is a group of translations with additional information.
type Fragment struct {
	Description  string
	Translations Translations
}

// FragmentMap presents fragementName => Fragment mapping.
type FragmentMap = map[string]Fragment

type ChineseConvertSetting struct {
	Mode string
	From string
	To   string
}

type ChineseSetting struct {
	Convert ChineseConvertSetting
}

type Settings struct {
	Chinese *ChineseSetting `yaml:",omitempty"`
}

type Output struct {
	Type string
}

// Outline is the struct of translation file.
type Outline struct {
	Settings  Settings
	Version   string
	Fragments FragmentMap
	Output    []Output
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
