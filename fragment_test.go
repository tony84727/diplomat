package diplomat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFragmentLint(t *testing.T) {
	fragment := Fragment{
		Translations: map[string]Translation{
			"key1": Translation{
				data: map[string]string{
					"zh-TW": "台灣",
					"en":    "Taiwan",
				},
			},
			"key2": Translation{
				data: map[string]string{
					"en": "WIP",
				},
			},
		},
	}
	errors := fragment.Lint()
	assert.Len(t, errors, 1)
}

func TestFragmentGetLocaleMap(t *testing.T) {
	fragment := Fragment{
		Translations: map[string]Translation{
			"admin": Translation{
				data: map[string]string{
					"zh-TW": "管理員",
					"en":    "Admin",
				},
			},
			"police": Translation{
				data: map[string]string{
					"zh-TW": "警察",
					"en":    "Police",
				},
			},
		},
	}
	localeMap := fragment.GetLocaleMap()
	chinese, exist := localeMap.Get("zh-TW")
	assert.True(t, exist)
	assert.Equal(t, "管理員", chinese.Translations["admin"])
	assert.Equal(t, "警察", chinese.Translations["police"])
	english, exist := localeMap.Get("en")
	assert.True(t, exist)
	assert.Equal(t, "Admin", english.Translations["admin"])
	assert.Equal(t, "Police", english.Translations["police"])
}
