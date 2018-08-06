package diplomat

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type TranslateEntry struct {
	Locale     string
	Translated string
}

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
func (t Translation) Iterate() <-chan TranslateEntry {
	c := make(chan TranslateEntry)
	m := make(map[string]string, len(t.data))
	for k, v := range t.data {
		m[k] = v
	}
	go func() {
		defer close(c)
		for k, v := range m {
			c <- TranslateEntry{k, v}
		}
	}()
	return c
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

type CopySetting struct {
	From string
	To   string
}

type Settings struct {
	Chinese *ChineseSetting `yaml:",omitempty"`
	Copy    []CopySetting
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
