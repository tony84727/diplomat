package yaml

import (
	"github.com/tony84727/diplomat/pkg/data"
)

type preprocessor struct {
	data.SimplePreprocessor
}

func (p *preprocessor) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct {
		Type    string      `yaml:"type"`
		Options interface{} `yaml:"options"`
	}
	if err := unmarshal(&actual); err != nil {
		return err
	}
	p.Type = actual.Type
	p.Options = actual.Options
	return nil
}

type templateOption struct {
	data.SimpleTemplateOption
}

func (t *templateOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var options map[string]interface{}
	if err := unmarshal(&options); err != nil {
		return err
	}
	t.SimpleTemplateOption = options
	return nil
}

type template struct {
	data.SimpleTemplate
}

func (t *template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct {
		Options templateOption `yaml:"options"`
		Type    string         `yaml:"type"`
	}
	if err := unmarshal(&actual); err != nil {
		return err
	}
	t.Type = actual.Type
	t.Options = actual.Options
	return nil
}

type output struct {
	data.SimpleOutput
}

func (o *output) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct {
		Selectors []string   `yaml:"selectors"`
		Templates []template `yaml:"templates"`
	}

	if err := unmarshal(&actual); err != nil {
		return err
	}
	o.Selectors = make([]data.Selector, len(actual.Selectors))
	for i, selector := range actual.Selectors {
		o.Selectors[i] = data.Selector(selector)
	}
	o.Templates = make([]data.Template, len(actual.Templates))
	for i, template := range actual.Templates {
		o.Templates[i] = template
	}
	return nil
}

type configurationFile struct {
	data.SimpleConfiguration
}

func (c *configurationFile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct {
		Outputs       []output       `yaml:"outputs"`
		Preprocessors []preprocessor `yaml:"preprocessors"`
	}
	if err := unmarshal(&actual); err != nil {
		return err
	}
	c.Outputs = make([]data.Output, len(actual.Outputs))
	for i, output := range actual.Outputs {
		c.Outputs[i] = output
	}
	c.Preprocessors = make([]data.Preprocessor, len(actual.Preprocessors))
	for i, preprocessor := range actual.Preprocessors {
		c.Preprocessors[i] = preprocessor
	}
	return nil
}
