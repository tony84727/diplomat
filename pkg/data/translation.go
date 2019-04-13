package data

type Translation interface {
	GetKey() string
	SetText(text string)
	GetText() *string
	GetChildren() []Translation
	AddChild(child Translation)
}
type translationNode struct {
	key string
	text *string
	children []Translation
}

func (t translationNode) GetKey() string {
	return t.key
}

func (t translationNode) GetText() *string {
	return t.text
}

func (t translationNode) GetChildren() []Translation {
	return t.children
}

func (t *translationNode) AddChild(child Translation) {
	t.children = append(t.children, child)
}

func (t *translationNode) SetText(text string) {
	t.text = &text
}

func NewTranslation(key string) Translation {
	return &translationNode{key:key,children:make([]Translation, 0)}
}