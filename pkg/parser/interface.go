package parser

import "github.com/insufficientchocolate/diplomat/pkg/data"

type Parser interface {
	GetTranslation() (*data.Translation,error)
}
