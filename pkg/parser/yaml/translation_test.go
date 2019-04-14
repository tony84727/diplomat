package yaml

import (
	"github.com/insufficientchocolate/diplomat/pkg/data"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type TranslationParserTestSuite struct {
	suite.Suite
}

func (p TranslationParserTestSuite) TestParse() {
	content, err := ioutil.ReadFile("testdata/admin.yaml")
	p.Require().NoError(err)
	parser := NewParser(content)
	translation, err := parser.GetTranslation()
	p.Require().NoError(err)
	walker := data.NewTranslationWalker(translation)
	p.ElementsMatch([]string{
		".admin.admin.zh-TW",
		".admin.admin.en",
		".admin.message.hello.zh-TW",
		".admin.message.hello.en",
	}, walker.GetKeys())
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, &TranslationParserTestSuite{})
}
