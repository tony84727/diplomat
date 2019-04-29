package types

import "github.com/tony84727/diplomat/pkg/data"

type Emitter interface {
	Emit(translation data.Translation, options data.TemplateOption) ([]byte, error)
}

type Factory interface {
	Build() Emitter
}
