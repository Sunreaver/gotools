package set

import (
	"bytes"
	"fmt"
	"sync"
)

// HashSet define HashSet
type HashSet struct {
	value map[interface{}]bool
	lock  sync.RWMutex
}

func newHashSet() *HashSet {
	return &HashSet{
		value: make(map[interface{}]bool),
		lock:  sync.RWMutex{},
	}
}

// Add will 添加，并返回是否已经添加
func (set *HashSet) Add(e interface{}) bool {
	set.lock.Lock()
	defer set.lock.Unlock()
	if r, ok := set.value[e]; !ok || !r {
		set.value[e] = true
		return true
	}
	return false
}

// AddArray will 添加一个序列，并返回真正添加的序列
func (set *HashSet) AddArray(e ...interface{}) []interface{} {
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

// Remove will 移除
func (set *HashSet) Remove(e interface{}) bool {
	set.lock.Lock()
	defer set.lock.Unlock()
	_, ok := set.value[e]
	delete(set.value, e)
	return ok
}

// Clear will 清空
func (set *HashSet) Clear() []interface{} {
	values := set.Elements()

	set.lock.Lock()
	defer set.lock.Unlock()
	set.value = make(map[interface{}]bool)
	return values
}

// Contains will 是否存在
func (set *HashSet) Contains(e interface{}) bool {
	set.lock.RLock()
	defer set.lock.RUnlock()
	_, ok := set.value[e]
	return ok
}

// Len will 返回是否存在
func (set *HashSet) Len() int {
	set.lock.RLock()
	defer set.lock.RUnlock()
	return len(set.value)
}

// Same will 判断是否相等
func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.value {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// Elements will 返回所有值
func (set *HashSet) Elements() []interface{} {
	set.lock.RLock()
	defer set.lock.RUnlock()

	initialLen := len(set.value)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key, value := range set.value {
		if !value {
			continue
		}
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

// String will 生成可视化输出
func (set *HashSet) String() string {
	set.lock.RLock()
	defer set.lock.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	for key := range set.value {
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
