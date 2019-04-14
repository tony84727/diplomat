package selector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombindSelector(t *testing.T) {
	a := NewPrefixSelector("admin", "hello")
	b := NewPrefixSelector("admin", "world")
	c := NewCombinedSelector(a, b)
	assert.True(t, c.IsValid([]string{"admin"}))
	assert.True(t, c.IsValid([]string{"admin", "hello"}))
	assert.True(t, c.IsValid([]string{"admin", "world"}))
	assert.False(t, c.IsValid([]string{"message", "hello"}))
}
