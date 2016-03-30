package set

import (
	"bytes"
	"fmt"
)

type HashSet map[interface{}]bool

func newHashSet() HashSet {
	return make(HashSet)
}

func (set *HashSet) Add(e interface{}) bool {
	if !(*set)[e] {
		(*set)[e] = true
		return true
	}
	return false
}

func (set *HashSet) AddArray(e []interface{}) []interface{} {
	index := 0
	result := make([]interface{}, len(e))
	for i := 0; i < len(e); i++ {
		ok := set.Add(e[i])
		if ok {
			result[index] = e[i]
			index++
		}
	}
	return result[:index]
}

func (set *HashSet) Remove(e interface{}) {
	delete((*set), e)
}

func (set *HashSet) Clear() {
	(*set) = make(HashSet)
}

func (set *HashSet) Contains(e interface{}) bool {
	return (*set)[e]
}

func (set *HashSet) Len() int {
	return len((*set))
}

func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if (*set).Len() != other.Len() {
		return false
	}
	for key := range *set {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (set *HashSet) Elements() []interface{} {
	initialLen := len((*set))
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range *set {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	for key := range *set {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}
