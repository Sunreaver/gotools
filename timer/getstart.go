package timer

import "time"

var (
	publicTimer *Timer
)

// CallBackType CallBack
type CallBackType func(e interface{})

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
		go publicTimer.AsyncStart()
	}
	publicTimer.Add(0, ts)

	return
}

// RemoveTimer 删除Timer
func RemoveTimer(id uint64) bool {
	deletedTimerSlice.Add(id)
	return true
}
