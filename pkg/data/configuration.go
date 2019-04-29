package data

type Preprocessor interface {
	GetType() string
	GetOptions() interface{}
}

type Selector string

type TemplateOption interface {
	GetFilename() string
	GetMapElement() map[string]interface{}
}

type Template interface {
	GetType() string
	GetOptions() TemplateOption
}

type Output interface {
	GetSelectors() []Selector
	GetTemplates() []Template
}

type Configuration interface {
	GetPreprocessors() []Preprocessor
	GetOutputs() []Output
}
