package parser

import "github.com/tony84727/diplomat/pkg/data"

type TranslationParser interface {
	GetTranslation() (data.Translation, error)
}

type ConfigurationParser interface {
	GetConfiguration() (data.Configuration, error)
}
