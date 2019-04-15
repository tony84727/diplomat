package transfrom

import (
	"github.com/insufficientchocolate/diplomat/pkg/data"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ChineseTransformerTestSuite struct {
	suite.Suite
}

func (c ChineseTransformerTestSuite) TestTransform() {
	transformer := NewChineseTransformer(ChineseTransformerOption{
		mode: TranditionalToSimplified,
		from: "zh-TW",
		to:"zh-CN",
	})
	translations := data.NewTranslation("")
	hello := data.NewTranslation("hello")
	translations.AddChild(hello)
	helloZhTW := data.NewTranslation("zh-TW")
	helloZhTW.SetText("學問")
	hello.AddChild(helloZhTW)
	c.Require().NoError(transformer.Transform(translations))
	exist := false
	walker := data.NewTranslationWalker(translations)
	walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
		if strings.Join(paths,".") == ".hello.zh-CN" {
			c.Equal("学问", *textNode.GetText())
			exist = true
		}
		return nil
	})
	c.True(exist)
}

func TestChineseTransformer(t *testing.T) {
	suite.Run(t, &ChineseTransformerTestSuite{})
}
