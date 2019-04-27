package selector

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
