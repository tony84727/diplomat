package chinese

import (
	"github.com/insufficientchocolate/diplomat/pkg/data"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type PreprocessorTestSuite struct {
	suite.Suite
}

func (p PreprocessorTestSuite) TestProcess() {
	instance := &Preprocessor{}
	root := data.NewTranslation("")
	questionNode := data.NewTranslation("question")
	root.AddChild(questionNode)
	traditional := data.NewTranslation("zh-TW")
	traditional.SetText("問題")
	questionNode.AddChild(traditional)
	p.NoError(instance.Process(root, map[interface{}]interface{}{
		"mode": "t2s",
		"from": "zh-TW",
		"to": "zh-CN",
	}))
	keys := make([]string, 0)
	walker := data.NewTranslationWalker(root)
	p.NoError(walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
		key :=  strings.Join(paths, ".")
		keys = append(keys,key)
		if key == "question.zh-CN" {
			p.Equal("问题", *textNode.GetText())
		}
		return nil
	}))
	p.ElementsMatch([]string{"question.zh-TW","question.zh-CN"}, keys)
}

func TestPreprocessor(t *testing.T) {
	suite.Run(t, &PreprocessorTestSuite{})
}
