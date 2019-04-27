package internal

import "github.com/tony84727/diplomat/pkg/data"

type PreprocessorFunc = func(translation data.Translation) error

func ComposePreprocessorFunc(funcs ...PreprocessorFunc) PreprocessorFunc {
	return func(translation data.Translation) error {
		for i := len(funcs) - 1; i >= 0; i-- {
			var err error
			err = funcs[i](translation)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
