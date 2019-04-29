package emit

import "github.com/tony84727/diplomat/pkg/emit/types"

var (
	GlobalRegistry Registry
)

type Registry interface {
	Get(name string) types.Emitter
	Registry(name string, instance types.Factory)
}

func NewRegistry() Registry {
	return &emitterRegistryImpl{
		instances: make(map[string]types.Emitter),
		factories: make(map[string]types.Factory),
	}
}

type emitterRegistryImpl struct {
	instances map[string]types.Emitter
	factories map[string]types.Factory
}

func (e *emitterRegistryImpl) Registry(name string, factory types.Factory) {
	e.factories[name] = factory
}

func (e *emitterRegistryImpl) Get(name string) types.Emitter {
	instance, exist := e.instances[name]
	if !exist {
		factory, exist := e.factories[name]
		if !exist {
			return nil
		}
		instance = factory.Build()
		e.instances[name] = instance
	}
	return instance
}

func init() {
	GlobalRegistry = NewRegistry()
}
