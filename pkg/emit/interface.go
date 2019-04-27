package emit

import "github.com/tony84727/diplomat/pkg/data"

type Emitter interface {
	Emit(translation data.Translation) ([]byte, error)
}
