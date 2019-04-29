package types

import "github.com/tony84727/diplomat/pkg/data"

type Emitter interface {
	Emit(translation data.Translation) ([]byte, error)
}

type Factory interface {
	Build() Emitter
}
