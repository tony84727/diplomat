package yaml

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/tony84727/diplomat/pkg/data"
	"io/ioutil"
	"testing"
)

type ConfigurationParserTestSuite struct {
	suite.Suite
}

func (c ConfigurationParserTestSuite) TestGetConfiguration() {
	content, err := ioutil.ReadFile("../../../testdata/diplomat.yaml")
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
	c.Equal(
		map[interface{}]interface{}{
			"mode": "t2s",
			"from": "zh-TW",
			"to":   "zh-CN",
		},
		preprocessors[0].GetOptions())
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
	c.Equal("js-object", template.GetType())
	c.Equal(templateOption{map[string]interface{}{"filename": "control-panel.js"}}, template.GetOptions())
}

func TestConfigurationParser(t *testing.T) {
	suite.Run(t, &ConfigurationParserTestSuite{})
}

type WriteTestSuite struct {
	suite.Suite
}

func (w WriteTestSuite) TestWrite() {
	content, err := ioutil.ReadFile("../../../testdata/diplomat.yaml")
	w.Require().NoError(err)
	parser := NewConfigurationParser(content)
	config, err := parser.GetConfiguration()
	w.Require().NoError(err)
	serialized, err := Write(config)
	w.Require().NoError(err)
	anotherParser := NewConfigurationParser(serialized)
	parsedConfig, err := anotherParser.GetConfiguration()
	w.Require().NoError(err)
	w.Equal(config, parsedConfig)
}

func TestWrite(t *testing.T) {
	suite.Run(t, &WriteTestSuite{})
}

func ExampleWrite() {
	content, err := ioutil.ReadFile("../../../testdata/diplomat.yaml")
	if err != nil {
		panic(err)
	}
	parser := NewConfigurationParser(content)
	config, err := parser.GetConfiguration()
	if err != nil {
		panic(err)
	}

	out, err := Write(config)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
	// Output:
	// preprocessors:
	// - type: chinese
	//   options:
	//     from: zh-TW
	//     mode: t2s
	//     to: zh-CN
	// - type: copy
	//   options:
	//     from: en
	//     to: fr
	// outputs:
	// - selectors:
	//   - admin
	//   - manage
	//   templates:
	//   - type: js-object
	//     options:
	//       filename: control-panel.js
}
