package yaml

import (
	"github.com/insufficientchocolate/diplomat/pkg/data"
)

type preprocessor struct {
	data.SimplePreprocessor
}

type templateOption struct {
	data.TemplateOption
}

type template struct {
	data.SimpleTemplate
}

func (t *template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct{
		options templateOption
		typeKey string `yaml:"type"`
	}
	if err := unmarshal(&actual); err != nil {
		return err
	}
	t.Type = actual.typeKey
	t.Options = actual.options
	return nil
}

type output struct {
	data.SimpleOutput
}

func (o *output) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct{
		selectors []string
		templates []template
	}

	if err := unmarshal(&actual);err != nil {
		return err
	}
	o.Selectors = make([]data.Selector, len(actual.selectors))
	for i, selector := range actual.selectors {
		o.Selectors[i] = data.Selector(selector)
	}
	o.Templates = make([]data.Template, len(actual.templates))
	for i, template := range actual.templates {
		o.Templates[i] = template
	}
	return nil
}

type configurationFile struct {
	data.SimpleConfiguration
}

func (c *configurationFile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var actual struct{
		outputs []output
		preprocessors []preprocessor
	}
	if err := unmarshal(&actual); err != nil {
		return err
	}
	c.Outputs = make([]data.Output,len(actual.outputs))
	for i, output := range actual.outputs {
		c.Outputs[i] = output
	}
	c.Preprocessors = make([]data.Preprocessor, len(actual.preprocessors))
	for i, preprocessor := range actual.preprocessors {
		c.Preprocessors[i] = preprocessor
	}
	return nil
}



