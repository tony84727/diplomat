package data

type TranslationWalker struct {
	root Translation
}

func NewTranslationWalker(root Translation) *TranslationWalker {
	return &TranslationWalker{root}
}

type worklistEntry struct {
	prefix *string
	root Translation
}

func (t TranslationWalker) GetKeys() []string {
	keys := make([]string, 0)
	worklist := []worklistEntry{{root:t.root}}
	for len(worklist) > 0 {
		start := worklist[0]
		worklist = worklist[1:]
		key := start.root.GetKey()
		var prefix string
		if start.prefix != nil {
			prefix = *start.prefix + "."
		}
		prefix = prefix + key
		if text := start.root.GetText(); text != nil {
			keys = append(keys, prefix)
		}
		for _, child := range start.root.GetChildren() {
			worklist = append(worklist, worklistEntry{prefix: &prefix,root:child})
		}
	}
	return keys
}
