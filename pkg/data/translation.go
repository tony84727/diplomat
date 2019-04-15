package data

type Translation interface {
	GetKey() string
	SetText(text string)
	GetText() *string
	GetChildren() map[string]Translation
	AddChild(child Translation)
	SetParent(parent Translation)
	GetParent() (parent Translation)
}
type translationNode struct {
	key string
	text *string
	children map[string]Translation
	parent Translation
}

func (t translationNode) GetKey() string {
	return t.key
}

func (t translationNode) GetText() *string {
	return t.text
}

func (t translationNode) GetChildren() map[string]Translation {
	return t.children
}

func (t *translationNode) AddChild(child Translation) {
	child.SetParent(t)
	t.children[child.GetKey()] = child
}

func (t *translationNode) SetText(text string) {
	t.text = &text
}

func (t *translationNode) SetParent(parent Translation) {
	t.parent = parent
}

func (t translationNode) GetParent() (parent Translation) {
	return t.parent
}

func NewTranslation(key string) Translation {
	return &translationNode{key:key,children:make(map[string]Translation)}
}