package data

type TranslationMerger struct {
	Translation
}

func (t TranslationMerger) Merge(other Translation) {
	walker := NewTranslationWalker(other)
	_ = walker.ForEachTextNode(func(paths []string, textNode Translation) error {
		var current Translation = t
		for _, segment := range paths[:len(paths)-1] {
			child := current.GetChild(segment)
			if child == nil {
				child = NewTranslation(segment)
				current.AddChild(child)
			}
			current = child
		}
		current.AddChild(textNode)
		return nil
	})
}

func NewTranslationMerger(root Translation) *TranslationMerger {
	return &TranslationMerger{root}
}

