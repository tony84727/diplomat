package selector

type PrefixSelector struct {
	keys []string
}

func (s PrefixSelector) IsValid(paths []string) bool {
	for i, s := range s.keys {
		if i >= len(paths) {
			break
		}
		if paths[i] != s {
			return false
		}
	}
	return true
}

func NewPrefixSelector(keys ...string) PrefixSelector {
	return PrefixSelector{
		keys,
	}
}
