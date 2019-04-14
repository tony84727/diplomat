package transfrom

import "github.com/insufficientchocolate/diplomat/pkg/data"

type Transformer interface {
	Transform(translation data.Translation) error
}
