package data

import "strings"

type Builder struct {
	Translation
}

func NewBuilder() *Builder {
	return &Builder{NewTranslation("")}
}

func (b Builder) Add(key string,text string) Translation {
	segments := strings.Split(key , ".")
	var current Translation = b
	for _, s := range segments {
		node := current.GetChild(s)
		if node == nil {
			node = NewTranslation(s)
			current.AddChild(node)
		}
		current = node
	}
	current.SetText(text)
	return current
}

