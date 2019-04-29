package types

import "github.com/tony84727/diplomat/pkg/data"

type Preprocessor interface {
	Process(translation data.Translation, option interface{}) error
}

type Factory interface {
	Build() Preprocessor
}
