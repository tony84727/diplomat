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
