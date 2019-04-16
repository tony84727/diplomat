package prepros


var (
	Manager *ManagerImpl
)

type ManagerImpl struct {
	preprocessors map[string]Preprocessor
}

func (m ManagerImpl) Registry(name string, instance Preprocessor) {
	m.preprocessors[name] = instance
}

func newManager() *ManagerImpl {
	return &ManagerImpl{
		preprocessors: make(map[string]Preprocessor),
	}
}

func init() {
	Manager = newManager()
}
