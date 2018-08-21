package diplomat

import (
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

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

func TestUnmarshalYAMLMap(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/admin.yaml")
	assert.NoError(t, err)
	yamlMap := make(YAMLMap)
	err = yaml.Unmarshal(data, yamlMap)
	assert.NoError(t, err)
	value, exist := yamlMap.GetKey("admin", "admin", "en")
	assert.True(t, exist)
	assert.Equal(t, "Admin", value)
	anotherMap, exist := yamlMap.GetKey("admin", "message", "hello")
	assert.True(t, exist)
	assert.IsType(t, yamlMap, anotherMap)
	value, exist = anotherMap.(YAMLMap).GetKey("zh-TW")
	assert.True(t, exist)
	assert.Equal(t, "您好", value)
}

func getSampleYAMLMap() YAMLMap {
	data := map[interface{}]interface{}{
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
	}
	return interfaceMapToYAMLMap(data)
}

func TestGetKeys(t *testing.T) {
	yamlMap := getSampleYAMLMap()
	assert.ElementsMatch(t, [][]string{
		[]string{"admin", "zh-TW"},
		[]string{"admin", "en"},
		[]string{"message", "hello", "zh-TW"},
		[]string{"message", "hello", "en"},
	}, yamlMap.GetKeys())
}

func TestFilterBySelector(t *testing.T) {
	yamlMap := getSampleYAMLMap()
	filtered := yamlMap.FilterBySelector(NewPrefixSelector("admin"))
	assert.Len(t, filtered, 1)
}

func TestHasKey(t *testing.T) {
	yamlMap := getSampleYAMLMap()
	assert.True(t, yamlMap.HasKey("admin", "zh-TW"))
	assert.True(t, yamlMap.HasKey("message", "hello"))
	assert.True(t, yamlMap.HasKey("message", "hello", "zh-TW"))
	assert.False(t, yamlMap.HasKey("hello", "admin"))
}
