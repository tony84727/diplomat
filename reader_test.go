package diplomat

import (
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOutline(t *testing.T) {
	outline := assertUnmarshalOutline(t, "testdata/diplomat.yaml")
	assert.Len(t, outline.Preprocessors, 2)
	assert.Len(t, outline.Output, 1)
}

func assertUnmarshalOutline(t *testing.T, path string) Outline {
	data, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	assert.NoError(t, err)
	return outline
}

func TestUnmarshalOutput(t *testing.T) {
	outline := assertUnmarshalOutline(t, "testdata/diplomat.yaml")
	assert.Len(t, outline.Output[0].Selectors, 2)
	filename, err := outline.Output[0].Template.Options.Get("filename")
	assert.NoError(t, err)
	assert.Equal(t, "control-panel.{{.Lang}}.js", filename)
}

func TestUnmarshalPreprocessorConfig(t *testing.T) {
	outline := assertUnmarshalOutline(t, "testdata/diplomat.yaml")
	assert.Equal(t, "chinese", outline.Preprocessors[0].Type)
	mode, err := outline.Preprocessors[0].Options.Get(0, "mode")
	assert.NoError(t, err)
	assert.Equal(t, "t2s", mode)
	from, err := outline.Preprocessors[0].Options.Get(0, "from")
	assert.NoError(t, err)
	assert.Equal(t, "zh-TW", from)
	to, err := outline.Preprocessors[0].Options.Get(0, "to")
	assert.NoError(t, err)
	assert.Equal(t, "zh-CN", to)
}
