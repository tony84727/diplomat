package emit

var (
	Registry EmitterRegistry
)

type EmitterRegistry struct {
	instances map[string]Emitter
}

func (e *EmitterRegistry) Registry(name string, instance Emitter) {
	e.instances[name] = instance
}

func (e *EmitterRegistry) Get(name string) Emitter {
	return e.instances[name]
}

func init() {
	Registry = EmitterRegistry{instances: make(map[string]Emitter)}
}

