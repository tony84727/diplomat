package prepros

import (
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/prepros/internal"
)

type Factory interface {
	Build() internal.PreprocessorFunc
}

type factory struct {
	preprocessor internal.PreprocessorFunc
}

func (f factory) Build() internal.PreprocessorFunc {
	return f.preprocessor
}

func NewFactory(registry Registry, configs ...data.Preprocessor) Factory {
	preprocessorInstances := make([]internal.PreprocessorFunc, 0, len(configs))
	// reverse order
	for i := len(configs) - 1; i >= 0; i-- {
		p := configs[i]
		if instance := registry.Get(p.GetType()); instance != nil {
			preprocessorInstances = append(preprocessorInstances, func(translation data.Translation) error {
				return instance.Process(translation, p.GetOptions())
			})
		}
	}
	return &factory{
		internal.ComposePreprocessorFunc(preprocessorInstances...),
	}
}
