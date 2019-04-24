package data

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type BuilderTestSuite struct {
	suite.Suite
}

func (b BuilderTestSuite) TestAdd() {
	builder := NewBuilder()
	translations := map[string]string{
		"hello.en":                           "Hello",
		"hello.zh-TW":                        "哈囉",
		"email.message.reset_password.en":    "Reset Password",
		"email.message.reset_password.zh-TW": "重設密碼",
	}
	for key, text := range translations {
		builder.Add(key, text)
	}

	walker := NewTranslationWalker(builder)
	collectedTranslation := make(map[string]string)
	_ = walker.ForEachTextNode(func(paths []string, textNode Translation) error {
		collectedTranslation[strings.Join(paths, ".")] = *textNode.GetText()
		return nil
	})
	b.Equal(translations, collectedTranslation)
}

func TestBuilder(t *testing.T) {
	suite.Run(t, &BuilderTestSuite{})
}
