package timer

import (
	"math"
	"sync"
	"time"
)

// dtsSize define DeleteTimerSlice Size,必须是2的N次方
const dtsSize uint64 = 0x10

var (
	publicTimer       *Timer
	deletedTimerSlice DeleteTimerSlice
)

func init() {
	deletedTimerSlice = make(DeleteTimerSlice, dtsSize)
	for i := 0; i < len(deletedTimerSlice); i++ {
		deletedTimerSlice[i].data = make(map[uint64]int64, 2048)
	}
}

// DeleteTimer 已经被删除的时间片
type DeleteTimer struct {
	sync.RWMutex
	data map[uint64]int64
}

// DeleteTimerSlice DeleteTimerSlice
type DeleteTimerSlice []DeleteTimer

// LastTime LastTime
func (d DeleteTimerSlice) LastTime(key uint64) int64 {
	index := key & (dtsSize - 1)
	d[index].RLock()
	defer d[index].RUnlock()
	if item, ok := d[index].data[key]; ok {
		return item
	}
	return math.MinInt64
}

// Add Add
func (d DeleteTimerSlice) Add(key uint64) {
	index := key & (dtsSize - 1)
	d[index].Lock()
	defer d[index].Unlock()
	d[index].data[key] = time.Now().UnixNano()
}

// Delete Delete
func (d DeleteTimerSlice) Delete(key uint64) {
	index := key & (dtsSize - 1)
	d[index].Lock()
	defer d[index].Unlock()
	delete(d[index].data, key)
}

// CallBackType CallBack
type CallBackType func(e interface{})

// Timer Timer
type Timer struct {
	Second TimeWheel
	Minute TimeWheel
	Hour   TimeWheel
}

// TimerSlice TimerSlice
type TimerSlice struct {
	CallBack     CallBackType
	Second       uint
	SecondOffset uint
	MinuteOffset uint
	Repeat       bool
	id           uint64
	insertTime   int64
	e            interface{}
}

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

func newTimer() *Timer {
	t := &Timer{
		Second: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Minute: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Hour:   TimeWheel{TS: make([][]*TimerSlice, 24), Index: 0},
	}

	go func(refTimer *Timer) {
		// 时间轮询
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 走一秒
				refTimer.Second.Tick()
				if refTimer.Second.Index == 0 {
					// 走一分钟
					refTimer.Minute.Tick()
					tmp := refTimer.Minute.CurTimerSliceAndClear()
					refTimer.Add(60, tmp...)
					if refTimer.Minute.Index == 0 {
						// 走一小时
						refTimer.Hour.Tick()
						tmp := refTimer.Hour.CurTimerSliceAndClear()
						refTimer.Add(3600, tmp...)
					}
				}

				//处理当前秒
				tmp := refTimer.Second.CurTimerSliceAndClear()
				go func(t *Timer, items []*TimerSlice) {
					for i := 0; i < len(items); i++ {
						item := items[i]
						if last := deletedTimerSlice.LastTime(item.id); last >= item.insertTime {
							// 此id插入时间比删除时间小,已经被删除
							continue
						}
						go func(t *Timer, instance *TimerSlice) {
							instance.CallBack(instance.e)
							if instance.Repeat {
								t.Add(0, instance)
							}
						}(t, item)
					}
				}(refTimer, tmp)

			}
		}
	}(t)
	return t
}

// Add 添加时间点到时间轮中
func (t *Timer) Add(seed uint, ts ...*TimerSlice) {
	for _, item := range ts {
		if last := deletedTimerSlice.LastTime(item.id); last >= item.insertTime {
			// 此id插入时间比删除时间小
			continue
		}

		second := item.Second
		if seed > 0 {
			if second < 60 {
				second %= seed
			} else if second < 3600 {
				second = (second + item.SecondOffset) % seed
			} else if second < 86400 {
				second = (second + item.SecondOffset + item.MinuteOffset*60) % seed
			}
		} else {
			item.SecondOffset = publicTimer.Second.Index
			item.MinuteOffset = publicTimer.Minute.Index
		}

		if second < 60 {
			//插入到秒轮
			index := (t.Second.Index + second) % 60
			t.Second.AddToIndex(uint(index), item)
		} else if second < 3600 {
			//插入到分钟轮
			index := (t.Minute.Index + second/60) % 60
			t.Minute.AddToIndex(uint(index), item)
		} else if second < 86400 {
			//插入到小时轮
			index := (t.Hour.Index + second/3600) % 24
			t.Hour.AddToIndex(uint(index), item)
		}
	}
}

// PutTimer 启动Timer
func PutTimer(second uint, repeat bool, id uint64, e interface{}, callBack CallBackType) {
	RemoveTimer(id)
	ts := &TimerSlice{
		CallBack:     callBack,
		Second:       second,
		SecondOffset: 0,
		MinuteOffset: 0,
		Repeat:       repeat,
		id:           id,
		insertTime:   time.Now().UnixNano(),
		e:            e,
	}
	if publicTimer == nil {
		publicTimer = newTimer()
	}
	publicTimer.Add(0, ts)

	return
}

// RemoveTimer 删除Timer
func RemoveTimer(id uint64) bool {
	deletedTimerSlice.Add(id)
	return true
}
