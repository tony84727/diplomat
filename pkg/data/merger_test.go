package data

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TranslationMergerTestSuite struct {
	suite.Suite
}

func (t TranslationMergerTestSuite) TestMerge() {
	root := NewTranslation("")
	hello := NewTranslation("hello")
	helloZhTW := NewTranslation("zh-TW")
	helloZhTW.SetText("您好")
	hello.AddChild(helloZhTW)
	root.AddChild(hello)

	anotherRoot := NewTranslation("")
	anotherHello := NewTranslation("hello")
	helloEN := NewTranslation("en")
	helloEN.SetText("hello")
	anotherHello.AddChild(helloEN)
	anotherRoot.AddChild(anotherHello)

	merger := NewTranslationMerger(root)
	merger.Merge(anotherRoot)

	walker := NewTranslationWalker(merger)
	keys := make([]string, 0)
	_ = walker.ForEachTextNode(func(path []string, textNode Translation) error {
		keys = append(keys, strings.Join(path, "."))
		return nil
	})
	t.ElementsMatch([]string{
		"hello.zh-TW",
		"hello.en",
	}, keys)
}

func TestTranslationMerger(t *testing.T) {
	suite.Run(t, &TranslationMergerTestSuite{})
}
