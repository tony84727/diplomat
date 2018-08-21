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

func TestYAMLOptionGet(t *testing.T) {
	option := YAMLOption{
		data: map[string]interface{}{
			"key1": []interface{}{1, 2, 3},
			"key2": map[string]interface{}{
				"a": "av",
				"b": "bv",
			},
		},
	}
	c, err := option.Get("key1", 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, c)
}

func TestYAMLOptionUnmarshal(t *testing.T) {
	content := `a: av
b: bv
c:
  - element1
  - element2
  - element3`
	var o YAMLOption
	err := yaml.Unmarshal([]byte(content), &o)
	assert.NoError(t, err)
	slice, err := o.Get("c")
	assert.NoError(t, err)
	assert.Len(t, slice, 3)
	assert.Equal(t, "element1", slice.([]interface{})[0])
	assert.Equal(t, "element2", slice.([]interface{})[1])
	assert.Equal(t, "element3", slice.([]interface{})[2])
	element3, err := o.Get("c", 2)
	assert.NoError(t, err)
	assert.Equal(t, "element3", element3)
}

// func TestYAMLOptionMarshal(t *testing.T) {
// 	option := YAMLOption{
// 		data: map[string]interface{}{
// 			"key1": []interface{}{1, 2, 3},
// 			"key2": map[string]interface{}{
// 				"a": "av",
// 				"b": "bv",
// 			},
// 		},
// 	}
// 	data, err := yaml.Marshal(option)
// 	assert.NoError(t, err)
// 	var o YAMLOption
// 	err = yaml.Unmarshal(data, &o)
// 	assert.NoError(t, err)
// 	assert.Equal(t, option, o)
// }

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

func getSampleNKV() NestedKeyValue {
	data, _ := nkvDataFromStringMap(map[string]interface{}{
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
	return NestedKeyValue{
		data: data,
	}
}

func TestGetKeys(t *testing.T) {
	nkv := getSampleNKV()
	assert.ElementsMatch(t, [][]string{
		[]string{"admin", "zh-TW"},
		[]string{"admin", "en"},
		[]string{"message", "hello", "zh-TW"},
		[]string{"message", "hello", "en"},
	}, nkv.GetKeys())
}

func TestFilterBySelector(t *testing.T) {
	nkv := getSampleNKV()
	filtered := nkv.FilterBySelector(NewPrefixSelector("admin"))
	assert.Len(t, filtered.data, 1)
}

func TestHasKey(t *testing.T) {
	nkv := getSampleNKV()
	assert.True(t, nkv.HasKey("admin", "zh-TW"))
	assert.True(t, nkv.HasKey("message", "hello"))
	assert.True(t, nkv.HasKey("message", "hello", "zh-TW"))
	assert.False(t, nkv.HasKey("hello", "admin"))
}
