package timer

import (
	"time"
)

// Timer Timer
type Timer struct {
	Second TimeWheel
	Minute TimeWheel
	Hour   TimeWheel
}

func newTimer() *Timer {
	t := &Timer{
		Second: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Minute: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Hour:   TimeWheel{TS: make([][]*TimerSlice, 24), Index: 0},
	}
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

// Tick will 时间轮滴答
func (t *Timer) Tick() {
	t.Second.Tick()
	if t.Second.Index == 0 {
		// 走一分钟
		t.Minute.Tick()
		tmp := t.Minute.CurTimerSliceAndClear()
		t.Add(60, tmp...)
		if t.Minute.Index == 0 {
			// 走一小时
			t.Hour.Tick()
			tmp := t.Hour.CurTimerSliceAndClear()
			t.Add(3600, tmp...)
		}
	}
}

// DoCurrent will 处理当前片
func (t *Timer) DoCurrent() {
	tmp := t.Second.CurTimerSliceAndClear()
	go func(t0 *Timer, items []*TimerSlice) {
		for i := 0; i < len(items); i++ {
			item := items[i]
			if last := deletedTimerSlice.LastTime(item.id); last >= item.insertTime {
				// 此id插入时间比删除时间小,已经被删除
				continue
			}
			go func(t1 *Timer, instance *TimerSlice) {
				// 实际执行
				instance.CallBack(instance.e)
				if instance.Repeat {
					t1.Add(0, instance)
				}
			}(t0, item)
		}
	}(t, tmp)
}

// AsyncStart will 同步方式启动时间轮
func (t *Timer) AsyncStart() {
	refTimer := t

	// 时间轮询
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 走一秒
			refTimer.Tick()
			// 处理当前秒
			refTimer.DoCurrent()
		}
	}
}
