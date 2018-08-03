package diplomat

import (
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

func TestTranslationUnmarshalYML(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/translation.yaml")
	assert.NoError(t, err)
	var translation Translation
	err = yaml.Unmarshal(data, &translation)
	assert.NoError(t, err)
	chinese, exist := translation.Get("zh-TW")
	assert.True(t, exist, "should have zh-TW translation")
	assert.Equal(t, "管理員", chinese)
	english, exist := translation.Get("en-US")
	assert.True(t, exist, "show have en-US translation")
	assert.Equal(t, "Admin", english)
}
