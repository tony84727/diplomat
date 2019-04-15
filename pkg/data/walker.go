package data

import "strings"

type TranslationWalker struct {
	root Translation
}

func NewTranslationWalker(root Translation) *TranslationWalker {
	return &TranslationWalker{root}
}

type worklistEntry struct {
	root Translation
	prefix []string
}

func (t TranslationWalker) GetKeys() []string {
	keys := make([]string,0)
	t.ForEachTextNode(func(paths []string,textNode Translation) error {
		keys = append(keys, strings.Join(paths, "."))
		return nil
	})
	return keys
}

func (t TranslationWalker) ForEachTextNode(callback func (paths []string,textNode Translation) error) error {
	worklist := []worklistEntry{{root:t.root, prefix: []string{}}}
	for len(worklist) > 0 {
		start := worklist[0]
		worklist = worklist[1:]
		key := start.root.GetKey()
		prefix := append(start.prefix,key)
		if text := start.root.GetText(); text != nil {
			if err := callback(prefix,start.root); err != nil {
				return err
			}
		}
		for _, child := range start.root.GetChildren() {
			worklist = append(worklist, worklistEntry{prefix: prefix,root:child})
		}
	}
	return nil
}
