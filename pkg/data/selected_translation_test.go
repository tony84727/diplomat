package data

import (
	"github.com/insufficientchocolate/diplomat/pkg/selector"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type SelectedTranslationTestSuite struct {
	suite.Suite
}

func (s SelectedTranslationTestSuite) TestWalk() {
	morning := map[string]string{
		"message.morning.greeting.zh-TW": "早安",
		"message.morning.greeting.en": "Good morning",
	}
	evening := map[string]string{
		"message.evening.greeting.zh-TW": "晚上好",
		"message.evening.greeting.en": "Good evening",
	}
	translationTree := NewBuilder()
	for key, text := range morning {
		translationTree.Add(key, text)
	}
	for key, text := range evening {
		translationTree.Add(key, text)
	}
	walker := NewTranslationWalker(NewSelectedTranslation(translationTree, selector.NewPrefixSelector("message","morning")))
	collected := make(map[string]string)
	_ = walker.ForEachTextNode(func(paths []string, textNode Translation) error {
		collected[strings.Join(paths,".")] = *textNode.GetText()
		return nil
	})
	s.Equal(morning, collected)
}

func TestSelectedTranslation(t *testing.T) {
	suite.Run(t, &SelectedTranslationTestSuite{})
}
