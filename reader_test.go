package diplomat

import (
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOutline(t *testing.T) {
	outline := assertUnmarshalOutline(t, "testdata/outline.yaml")
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
	outline := assertUnmarshalOutline(t, "testdata/outline.yaml")
	assert.Len(t, outline.Output[0].Selectors, 2)
	filename, err := outline.Output[0].Template.Options.Get("filename")
	assert.NoError(t, err)
	assert.Equal(t, "{{.Locale}}.{{.FragmentName}}.js", filename)
}

func TestUnmarshalPreprocessorConfig(t *testing.T) {
	outline := assertUnmarshalOutline(t, "testdata/outline.yaml")
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

func TestYAMLOptionGet(t *testing.T) {
	option := YAMLOption{
		data: map[interface{}]interface{}{
			"key1": []interface{}{1, 2, 3},
			"key2": map[interface{}]interface{}{
				"a": "av",
				"b": "bv",
			},
		},
	}
	c, err := option.Get("key1", 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, c)
}

func TestMarshalYAMLOption(t *testing.T) {
	option := YAMLOption{
		data: map[interface{}]interface{}{
			"key1": []interface{}{1, 2, 3},
			"key2": map[interface{}]interface{}{
				"a": "av",
				"b": "bv",
			},
		},
	}
	data, err := yaml.Marshal(option)
	assert.NoError(t, err)
	var out YAMLOption
	err = yaml.Unmarshal(data, &out)
	assert.NoError(t, err)
	assert.Equal(t, option, out)
}

func TestOutlineMarshalAndUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/outline.yaml")
	assert.NoError(t, err)
	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	assert.NoError(t, err)
	output, err := yaml.Marshal(outline)
	assert.NoError(t, err)
	assert.Equal(t, string(data), string(output))
}

func TestUnmarshalNestedKeyValue(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/admin.yaml")
	assert.NoError(t, err)
	var nkv NestedKeyValue
	err = yaml.Unmarshal(data, &nkv)
	assert.NoError(t, err)
	value, exist := nkv.GetKey("admin", "admin", "en")
	assert.True(t, exist)
	assert.Equal(t, "Admin", value)
	anotherMap, exist := nkv.GetKey("admin", "message", "hello")
	assert.True(t, exist)
	assert.IsType(t, nkv, anotherMap)
	value, exist = anotherMap.(NestedKeyValue).GetKey("zh-TW")
	assert.True(t, exist)
	assert.Equal(t, "您好", value)
}

func TestGetKeys(t *testing.T) {
	data, err := nkvDataFromStringMap(map[string]interface{}{
		"admin": map[interface{}]interface{}{
			"zh-TW": "管理員",
			"en":    "admin",
		},
		"message": map[interface{}]interface{}{
			"hello": map[interface{}]interface{}{
				"zh-TW": "您好",
				"en":    "Hello!",
			},
		},
	})
	assert.NoError(t, err)
	nkv := NestedKeyValue{
		data: data,
	}
	assert.ElementsMatch(t, [][]string{
		[]string{"admin", "zh-TW"},
		[]string{"admin", "en"},
		[]string{"message", "hello", "zh-TW"},
		[]string{"message", "hello", "en"},
	}, nkv.GetKeys())
}
