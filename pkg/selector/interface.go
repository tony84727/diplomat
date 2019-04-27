package selector

type Selector interface {
	IsValid(paths []string) bool
}
