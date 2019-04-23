package emit

var (
	GlobalRegistry Registry
)

type Registry interface {
	Get(name string) Emitter
	Registry(name string,instance Emitter)
}

type emitterRegistryImpl struct {
	instances map[string]Emitter
}

func (e *emitterRegistryImpl) Registry(name string, instance Emitter) {
	e.instances[name] = instance
}

func (e *emitterRegistryImpl) Get(name string) Emitter {
	return e.instances[name]
}

func init() {
	GlobalRegistry = &emitterRegistryImpl{instances: make(map[string]Emitter)}
}

