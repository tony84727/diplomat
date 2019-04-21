package data

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type translationWalkerTestSuite struct {
	suite.Suite
}

func (t translationWalkerTestSuite) TestGetKeys() {
	root := NewTranslation("")
	hello := NewTranslation("hello")
	helloEnglish := NewTranslation("english")
	helloEnglish.SetText("hello")
	chineseEnglish := NewTranslation("chinese")
	chineseEnglish.SetText("您好")
	hello.AddChild(helloEnglish)
	hello.AddChild(chineseEnglish)
	root.AddChild(hello)
	walker := NewTranslationWalker(root)
	keys := walker.GetKeys()
	t.ElementsMatch([]string{"hello.english","hello.chinese"},keys)
}

func (t translationWalkerTestSuite) TestBacktracking() {
	root := NewTranslation("")
	hello := NewTranslation("hello")
	helloEnglish := NewTranslation("en")
	helloEnglish.SetText("Hello")
	hello.AddChild(helloEnglish)
	root.AddChild(hello)
	world := NewTranslation("world")
	worldEnglish := NewTranslation("en")
	worldEnglish.SetText("World")
	world.AddChild(worldEnglish)
	root.AddChild(world)
	walker := NewTranslationWalker(root)
	keys := make([]string, 0)
	t.NoError(walker.ForEachTextNodeWithBacktracking(func(paths []string, textNode Translation) error {
		keys = append(keys, strings.Join(paths, "."))
		return nil
	}, func(paths []string) bool {
		return paths[0] != "world"
	}))
	t.ElementsMatch([]string{"hello.en"}, keys)
}

func TestTranslationWalker(t *testing.T) {
	suite.Run(t, &translationWalkerTestSuite{})
}
