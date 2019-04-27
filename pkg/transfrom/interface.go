package transfrom

import "github.com/tony84727/diplomat/pkg/data"

type Transformer interface {
	Transform(translation data.Translation) error
}
