package selector

import "github.com/insufficientchocolate/diplomat/pkg/data"

type selectedTranslation struct {
	Selector
	data.Translation
	shallowTree data.Translation
}


func NewSelectedTranslation(origin data.Translation, selector Selector) data.Translation {
	shallowTree := data.NewTranslation("")
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
