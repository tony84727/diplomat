package types

type FactoryWrapper struct {
	Constructor func() Preprocessor
}

func (f FactoryWrapper) Build() Preprocessor {
	return f.Constructor()
}
