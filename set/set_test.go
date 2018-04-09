package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	a Set
	b Set
	c Set
)

func TestMain(t *testing.T) {
	a = NewSimpleSet()
	a.AddArray("a", "b", "c", 1)
	b = NewSimpleSet()
	b.AddArray("c", "d")
	c = NewSimpleSet()
	c.AddArray("a", 1)
}

func TestIsSuperset(t *testing.T) {
	assert.Equal(t, false, IsSuperset(a, b))
	assert.Equal(t, true, IsSuperset(a, c))
	assert.Equal(t, false, IsSuperset(a, a))
}

func TestUnion(t *testing.T) {
	d := Union(a, b)
	assert.Equal(t, true, d.Contains("a"))
	assert.Equal(t, true, d.Contains("b"))
	assert.Equal(t, true, d.Contains("c"))
	assert.Equal(t, true, d.Contains("d"))
	assert.Equal(t, true, d.Contains(1))
	assert.Equal(t, false, d.Contains("e"))
}

func TestIntersect(t *testing.T) {
	d := Intersect(a, b)
	assert.Equal(t, false, d.Contains("a"))
	assert.Equal(t, false, d.Contains("b"))
	assert.Equal(t, true, d.Contains("c"))
	assert.Equal(t, false, d.Contains("d"))
	assert.Equal(t, false, d.Contains(1))
	assert.Equal(t, false, d.Contains("e"))
}

func TestDifference(t *testing.T) {
	d := Difference(a, b)
	assert.Equal(t, true, d.Contains("a"))
	assert.Equal(t, true, d.Contains("b"))
	assert.Equal(t, false, d.Contains("c"))
	assert.Equal(t, false, d.Contains("d"))
	assert.Equal(t, true, d.Contains(1))
	assert.Equal(t, false, d.Contains("e"))
}

func TestSymmetricDifference(t *testing.T) {
	d := SymmetricDifference(a, b)
	assert.Equal(t, true, d.Contains("a"))
	assert.Equal(t, true, d.Contains("b"))
	assert.Equal(t, false, d.Contains("c"))
	assert.Equal(t, true, d.Contains("d"))
	assert.Equal(t, true, d.Contains(1))
	assert.Equal(t, false, d.Contains("e"))
}

func TestIsSet(t *testing.T) {
	assert.Equal(t, true, IsSet(a))
	assert.Equal(t, true, IsSet(b))
	assert.Equal(t, false, IsSet(2))
	assert.Equal(t, false, IsSet([]Set{a, c}))
}
