package types

type FactoryWrapper struct {
	Constructor func() Emitter
}

func (f FactoryWrapper) Build() Emitter {
	return f.Constructor()
}
