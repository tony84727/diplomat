package data

type TranslationMerger struct {
	Translation
}

func (t TranslationMerger) Merge(other Translation) {
	panic("implement me")
}

func NewTranslationMerger(root Translation) *TranslationMerger {
	return &TranslationMerger{root}
}

