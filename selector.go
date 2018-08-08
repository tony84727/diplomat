package diplomat

type Selector interface {
	IsValid(paths []string) bool
}

type PrefixSelector struct {
	keys []string
}

func (s PrefixSelector) IsValid(paths []string) bool {
	for i, s := range s.keys {
		if paths[i] != s {
			return false
		}
	}
	return true
}

type CombinedSelector struct {
	selectors []Selector
}

func (c CombinedSelector) IsValid(paths []string) bool {
	for _, s := range c.selectors {
		if s.IsValid(paths) {
			return true
		}
	}
	return false
}

func NewCombinedSelector(selectors ...Selector) Selector {
	return CombinedSelector{
		selectors,
	}
}
