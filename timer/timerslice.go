package timer

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
