package data

type SimpleConfiguration struct {
	Preprocessors []Preprocessor
	Outputs []Output
}

func (s SimpleConfiguration) GetPreprocessors() []Preprocessor {
	return s.Preprocessors
}

func (s SimpleConfiguration) GetOutputs() []Output {
	return s.Outputs
}

type SimpleOutput struct {
	Selectors []Selector
	Templates []Template
}

func (s SimpleOutput) GetSelectors() []Selector {
	return s.Selectors
}

func (s SimpleOutput) GetTemplates() []Template {
	return s.Templates
}

type SimpleTemplate struct {
	Type string
	Options TemplateOption
}

func (s SimpleTemplate) GetType() string {
	return s.Type
}

func (s SimpleTemplate) GetOptions() TemplateOption {
	return s.Options
}

type SimpleTemplateOption string

func (s SimpleTemplateOption) GetFilename() string {
	return string(s)
}

type SimplePreprocessor struct {
	Type string
	Options interface{}
}

func (s SimplePreprocessor) GetType() string {
	return s.Type
}

func (s SimplePreprocessor) GetOptions() interface{} {
	return s.Options
}



