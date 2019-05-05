package internal

import (
	"github.com/stretchr/testify/suite"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
	"io/ioutil"
	"testing"
)

type ConfigNavigatorTestSuite struct {
	suite.Suite
}

func (c ConfigNavigatorTestSuite) TestGet() {
	testConfigContent, err := ioutil.ReadFile("../../../testdata/diplomat.yaml")
	c.Require().NoError(err)
	paser := yaml.NewConfigurationParser(testConfigContent)
	config, err := paser.GetConfiguration()
	c.Require().NoError(err)
	navigator := NewConfigNavigator(config)
	mode, err := navigator.Get("preprocessors","chinese","options","mode")
	c.Require().NoError(err)
	c.Equal("t2s", mode)
}

func TestConfigNavigator(t *testing.T) {
	suite.Run(t, &ConfigNavigatorTestSuite{})
}
