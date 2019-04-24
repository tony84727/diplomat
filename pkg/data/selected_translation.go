package data

import (
	"github.com/tony84727/diplomat/pkg/selector"
	"strings"
)

type selectedTranslation struct {
	selector.Selector
	Translation
	shallowTree Translation
}

func NewSelectedTranslation(origin Translation, selector selector.Selector) Translation {
	walker := NewTranslationWalker(origin)
	root := NewBuilder()
	_ = walker.ForEachTextNodeWithBacktracking(func(paths []string, textNode Translation) error {
		root.Add(strings.Join(paths, "."), *textNode.GetText())
		return nil
	}, func(paths []string) bool {
		return selector.IsValid(paths)
	})
	return root
}
