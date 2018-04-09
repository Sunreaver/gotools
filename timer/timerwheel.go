package timer

import "sync"

// TimeWheel 时间轮
type TimeWheel struct {
	sync.RWMutex
	TS    [][]*TimerSlice
	Index uint
}

// AddToIndex 给当前时间轮添加点
func (tw *TimeWheel) AddToIndex(index uint, ts *TimerSlice) bool {
	tw.Lock()
	defer tw.Unlock()
	if index > uint(len(tw.TS)) {
		return false
	}
	tw.TS[index] = append(tw.TS[index], ts)
	return true
}

func (tw *TimeWheel) unsafeRemove(index, at uint) bool {
	if index > uint(len(tw.TS)) || at > uint(len(tw.TS[index])) {
		return false
	}

	tw.TS[index] = append(tw.TS[index][:at], tw.TS[index][at+1:]...)
	return true
}

// RemoveWithID 删除某个ID
func (tw *TimeWheel) RemoveWithID(id uint64) bool {
	tw.Lock()
	defer tw.Unlock()
	for index := 0; index < len(tw.TS); index++ {
		t := tw.TS[index]
		for i := 0; i < len(t); i++ {
			if t[i].id == id {
				tw.unsafeRemove(uint(index), uint(i))
				return true
			}
		}
	}
	return false
}

// CurTimerSliceAndClear 获取当前时间片，且清空时间片
func (tw *TimeWheel) CurTimerSliceAndClear() []*TimerSlice {
	tw.Lock()
	tmp := tw.TS[tw.Index]
	tw.TS[tw.Index] = []*TimerSlice{}
	tw.Unlock()
	return tmp
}

// Tick 时间轮自增一次
func (tw *TimeWheel) Tick() {
	tw.Lock()
	tw.Index = (tw.Index + 1) % uint(len(tw.TS))
	tw.Unlock()
}
