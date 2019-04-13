package data

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type translationWalkerTestSuite struct {
	suite.Suite
}

func (t translationWalkerTestSuite) TestGetKeys() {
	root := NewTranslation("root")
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
	t.ElementsMatch([]string{"root.hello.english","root.hello.chinese"},keys)
}

func TestTranslationWalker(t *testing.T) {
	suite.Run(t, &translationWalkerTestSuite{})
}
