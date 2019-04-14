package parser

import "github.com/insufficientchocolate/diplomat/pkg/data"

type TranslationParser interface {
	GetTranslation() (data.Translation,error)
}
