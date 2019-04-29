package prepros

import "github.com/tony84727/diplomat/pkg/prepros/types"

var (
	GlobalRegistry Registry
)

type Registry interface {
	Get(name string) types.Preprocessor
	Registry(name string, instance types.Factory)
}

type registryImpl struct {
	instances map[string]types.Preprocessor
	factories map[string]types.Factory
}

func (m registryImpl) Registry(name string, factory types.Factory) {
	m.factories[name] = factory
}

func (m registryImpl) Get(name string) types.Preprocessor {
	instance, exist := m.instances[name]
	if !exist {
		factory, exist := m.factories[name]
		if !exist {
			return nil
		}
		instance = factory.Build()
		m.instances[name] = instance
	}
	return instance
}

func newRegistry() *registryImpl {
	return &registryImpl{
		instances: make(map[string]types.Preprocessor),
		factories: make(map[string]types.Factory),
	}
}

func init() {
	GlobalRegistry = newRegistry()
}
