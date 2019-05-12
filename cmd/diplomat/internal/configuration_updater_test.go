package internal

import (
	"github.com/stretchr/testify/suite"
	"github.com/tony84727/diplomat/pkg/data"
	"testing"
)

type ConfigurationUpdaterTestSuite struct {
	suite.Suite
}

func (c ConfigurationUpdaterTestSuite) TestSet() {
	config := &data.SimpleConfiguration{
		Preprocessors: []data.Preprocessor{
			&data.SimplePreprocessor{
				Type: "dummy",
				Options: map[string]string{
					"from": "zh-TW",
				},
			},
		},
		Outputs: make([]data.Output, 0),
	}
	updater := NewConfigurationUpdater(config)
	c.Require().NoError(updater.Set("preprocessors.0.options.from", "en"))
	c.Equal("en",config.Preprocessors[0].GetOptions().(map[string]string)["from"])
}

func TestConfigurationUpdater(t *testing.T) {
	suite.Run(t, &ConfigurationUpdaterTestSuite{})
}