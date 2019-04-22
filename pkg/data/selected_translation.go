package data

import "github.com/insufficientchocolate/diplomat/pkg/selector"

type selectedTranslation struct {
	selector.Selector
	Translation
	shallowTree Translation
}


func NewSelectedTranslation(origin Translation, selector selector.Selector) Translation {
	shallowTree := NewTranslation("")
	for _, c := range origin.GetChildren() {
		shallowTree.AddChild(c)
		// FIXME: AddChild will aslo set parent. Set it back to origin here. Maybe introduce a flag to disable auto setting parent
		c.SetParent(origin)
	}
	return &selectedTranslation{
		Selector: selector,
		Translation: origin,
		shallowTree: shallowTree,
	}
}
