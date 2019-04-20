package data

type Translation interface {
	GetKey() string
	SetText(text string)
	GetText() *string
	GetChildren() []Translation
	GetChild(key string) Translation
	AddChild(child Translation)
	SetParent(parent Translation)
	GetParent() (parent Translation)
}
type translationNode struct {
	key string
	text *string
	children []Translation
	keyIndex map[string]int
	parent Translation
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

func (t translationNode) GetChild(key string) Translation {
	index, exist := t.keyIndex[key]
	if !exist {
		return nil
	}
	return t.children[index]
}

func (t *translationNode) AddChild(child Translation) {
	child.SetParent(t)
	key := child.GetKey()
	index, exist := t.keyIndex[key]
	if exist {
		t.children[index] = child
	}
	t.children = append(t.children, child)
	t.keyIndex[key] = len(t.children) - 1
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
	return &translationNode{key:key,children:make([]Translation, 0), keyIndex: make(map[string]int)}
}