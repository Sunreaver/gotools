package set

// Set define Set
type Set interface {
	Add(e interface{}) (added bool)
	AddArray(e ...interface{}) (addedKeys []interface{})
	Remove(e interface{}) (exist bool)
	Clear() (keys []interface{})
	Contains(e interface{}) (exist bool)
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}

// IsSuperset 判断集合 one 是否是集合 other 的超集
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}

// Union 生成集合 one 和集合 other 的并集
func Union(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	unionedSet := NewSimpleSet()
	for _, v := range one.Elements() {
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range other.Elements() {
		unionedSet.Add(v)
	}
	return unionedSet
}

// Intersect 生成集合 one 和集合 other 的交集
func Intersect(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	intersectedSet := NewSimpleSet()
	if other.Len() == 0 {
		return intersectedSet
	}
	if one.Len() < other.Len() {
		for _, v := range one.Elements() {
			if other.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	} else {
		for _, v := range other.Elements() {
			if one.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}

// Difference 生成集合 one 对集合 other 的差集
func Difference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	differencedSet := NewSimpleSet()
	if other.Len() == 0 {
		for _, v := range one.Elements() {
			differencedSet.Add(v)
		}
		return differencedSet
	}
	for _, v := range one.Elements() {
		if !other.Contains(v) {
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

// SymmetricDifference 生成集合 one 和集合 other 的对称差集
func SymmetricDifference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	diffA := Difference(one, other)
	if other.Len() == 0 {
		return diffA
	}
	diffB := Difference(other, one)
	return Union(diffA, diffB)
}

// NewSimpleSet will 生成新set
func NewSimpleSet() Set {
	hash := newHashSet()
	return hash
}

// IsSet will 判断是否为Set
func IsSet(value interface{}) bool {
	if _, ok := value.(Set); ok {
		return true
	}
	return false
}
