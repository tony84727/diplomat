package diplomat

import (
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

func TestOptionToChineseTransformer(t *testing.T) {
	optionContent := `
mode: t2s
from: zh-TW
to: zh-CN`
	var o YAMLOption
	err := yaml.Unmarshal([]byte(optionContent), &o)
	assert.NoError(t, err)
	transformer, err := optionToChineseTransformer(o)
	assert.NoError(t, err)
	assert.Equal(t, TranditionalToSimplified, transformer.mode)
	assert.Equal(t, "zh-TW", transformer.from)
	assert.Equal(t, "zh-CN", transformer.to)
}

func TestChineseConvertorPreprocessorFunc(t *testing.T) {
	transformer := ChineseTransformer{
		from: "zh-TW",
		to:   "zh-CN",
		mode: TranditionalToSimplified,
	}
	handler := transformer.getPreprocessorFunc()
	nkv := getSampleNKV()
	err := handler(&nkv)
	assert.NoError(t, err)
	cn, exist := nkv.GetKey("admin", "zh-CN")
	assert.True(t, exist)
	assert.Equal(t, "管理员", cn)
}
