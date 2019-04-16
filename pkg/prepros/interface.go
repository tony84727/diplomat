package prepros

import "github.com/insufficientchocolate/diplomat/pkg/data"

type Preprocessor interface {
	Process(translation data.Translation, option interface{}) error
}
