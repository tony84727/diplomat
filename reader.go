package diplomat

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Translation presents languageCode => translatedText mapping
type Translation struct {
	data map[string]string
}

func (t Translation) Get(locale string) (translated string, exist bool) {
	translated, exist = t.data[locale]
	return
}

func (t Translation) Set(locale, translated string) {
	t.data[locale] = translated
}

func (t Translation) GetLocales() []string {
	locales := make([]string, 0, 1)
	for locale := range t.data {
		locales = append(locales, locale)
	}
	return locales
}

func (t *Translation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&t.data)
}

// Translations presents translationKey => translation (for different language) mapping
type Translations = map[string]Translation

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

type FragmentOutputSetting struct {
	Type string
	Name string
}

type OutputSetting struct {
	Fragments []FragmentOutputSetting
}

// Outline is the struct of translation file.
type Outline struct {
	Settings  Settings
	Version   string
	Fragments FragmentMap
	Output    OutputSetting
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
