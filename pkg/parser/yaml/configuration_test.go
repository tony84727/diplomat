package yaml

import (
	"github.com/stretchr/testify/suite"
	"github.com/tony84727/diplomat/pkg/data"
	"io/ioutil"
	"testing"
)

type ConfigurationParserTestSuite struct {
	suite.Suite
}

func (c ConfigurationParserTestSuite) TestGetConfiguration() {
	content, err := ioutil.ReadFile("testdata/diplomat.yaml")
	c.Require().NoError(err)
	parser := NewConfigurationParser(content)
	configuration, err := parser.GetConfiguration()
	c.Require().NoError(err)
	// preprocessors
	preprocessors := configuration.GetPreprocessors()
	preprocessorTypes := make([]string, len(preprocessors))
	for i, p := range preprocessors {
		preprocessorTypes[i] = p.GetType()
	}
	c.ElementsMatch([]string{"chinese", "copy"}, preprocessorTypes)
	c.Equal([]interface{}{
		map[interface{}]interface{}{
			"mode": "t2s",
			"from": "zh-TW",
			"to":   "zh-CN",
		},
	}, preprocessors[0].GetOptions())
	c.Equal(
		map[interface{}]interface{}{"from": "en", "to": "fr"},
		preprocessors[1].GetOptions())
	// outputs
	outputs := configuration.GetOutputs()
	c.Require().Len(outputs, 1)
	output := outputs[0]
	c.ElementsMatch([]data.Selector{"admin", "manage"}, output.GetSelectors())
	c.Require().Len(output.GetTemplates(), 1)
	template := outputs[0].GetTemplates()[0]
	c.Equal("js", template.GetType())
	c.Equal(templateOption{data.SimpleTemplateOption("control-panel.{{.Lang}}.js")}, template.GetOptions())
}

func TestConfigurationParser(t *testing.T) {
	suite.Run(t, &ConfigurationParserTestSuite{})
}
