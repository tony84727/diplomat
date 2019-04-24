package selector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrefixSelector(t *testing.T) {
	s := NewPrefixSelector("admin", "hello", "world")
	assert.True(t, s.IsValid([]string{"admin"}))
	assert.True(t, s.IsValid([]string{"admin", "hello"}))
	assert.True(t, s.IsValid([]string{"admin", "hello", "world"}))
	assert.False(t, s.IsValid([]string{"message"}))
	assert.False(t, s.IsValid([]string{"hello"}))
}
