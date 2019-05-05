package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
	"io/ioutil"
	"reflect"
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
	mode, err := navigator.Get("preprocessors","0","options","mode")
	c.Require().NoError(err)
	c.Equal("t2s", mode)
}

func TestConfigNavigator(t *testing.T) {
	suite.Run(t, &ConfigNavigatorTestSuite{})
}

func Test_searchField(t *testing.T) {
	fake := struct {
		Name string `navigate:"first"`
		Number int `navigate:"second"`
	}{
		"whatever",
		100,
	}
	i, ok := searchField("first", reflect.TypeOf(fake))
	require.True(t, ok)
	assert.Equal(t, []int{0},i)
}

func Test_searchField_Nested(t *testing.T) {
	type nested struct {
		Name string `navigate:"first"`
		Number int `navigate:"second"`
	}
	fake := struct {
		embedded nested
	}{
		embedded: nested{Name: "whatever", Number: 100 },
	}
	i, ok := searchField("second", reflect.TypeOf(fake))
	require.True(t, ok)
	assert.Equal(t, []int{0,1}, i)
}
