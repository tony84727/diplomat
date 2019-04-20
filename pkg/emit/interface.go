package emit

import "github.com/insufficientchocolate/diplomat/pkg/data"

type Emitter interface {
	Emit(translation data.Translation) ([]byte,error)
}
