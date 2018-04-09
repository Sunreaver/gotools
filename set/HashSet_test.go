package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	s := newHashSet()
	assert.Equal(t, true, s.Add("a"))
	assert.Equal(t, false, s.Add("a"))
}

func TestAddArray(t *testing.T) {
	s := newHashSet()
	assert.Equal(t, []interface{}{"a", 1}, s.AddArray("a", 1))
	assert.Equal(t, []interface{}{"b", "c"}, s.AddArray([]interface{}{"a", "b", 1, "c"}...))
}

func TestRemove(t *testing.T) {
	s := newHashSet()
	assert.Equal(t, []interface{}{"a", 1}, s.AddArray("a", 1))
	assert.Equal(t, false, s.Remove("b"))
	assert.Equal(t, true, s.Remove("a"))
}

func TestClear(t *testing.T) {
	s := newHashSet()
	assert.Equal(t, []interface{}{"a", 1}, s.AddArray("a", 1))
	assert.Equal(t, 2, len(s.Clear()))
	assert.Equal(t, []interface{}{}, s.Elements())
}
