package diplomat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixSelector(t *testing.T) {
	s := NewPrefixSelector("admin", "hello", "world")
	assert.True(t, s.IsValid([]string{"admin"}))
	assert.True(t, s.IsValid([]string{"admin", "hello"}))
	assert.True(t, s.IsValid([]string{"admin", "hello", "world"}))
	assert.False(t, s.IsValid([]string{"message"}))
	assert.False(t, s.IsValid([]string{"hello"}))
}

func TestCombindSelector(t *testing.T) {
	a := NewPrefixSelector("admin", "hello")
	b := NewPrefixSelector("admin", "world")
	c := NewCombinedSelector(a, b)
	assert.True(t, c.IsValid([]string{"admin"}))
	assert.True(t, c.IsValid([]string{"admin", "hello"}))
	assert.True(t, c.IsValid([]string{"admin", "world"}))
	assert.False(t, c.IsValid([]string{"message", "hello"}))
}
