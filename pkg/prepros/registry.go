package prepros

var (
	GlobalRegistry Registry
)

type Registry interface {
	Get(name string) Preprocessor
	Registry(name string, instance Preprocessor)
}

type managerImpl struct {
	preprocessors map[string]Preprocessor
}

func (m managerImpl) Registry(name string, instance Preprocessor) {
	m.preprocessors[name] = instance
}

func (m managerImpl) Get(name string) Preprocessor {
	return m.preprocessors[name]
}

func newManager() *managerImpl {
	return &managerImpl{
		preprocessors: make(map[string]Preprocessor),
	}
}

func init() {
	GlobalRegistry = newManager()
}
