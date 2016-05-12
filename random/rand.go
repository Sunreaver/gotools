package random

import "time"

const (
	_MTLen = 624
)

var (
	_MT   []int
	index = 0
)

func init() {
	_MT = make([]int, _MTLen)
	_MT[0] = int(time.Now().Unix())
	for i := 1; i < _MTLen; i++ {
		// MT[i] := last 32 bits of(1812433253 * (MT[i-1] xor (right shift by 30 bits(MT[i-1]))) + i) // 0x6c078965
		_MT[i] = int((int64(0x6c078965)*(int64(_MT[i-1])^int64(_MT[i-1]>>30)) + int64(i)) & 0x00000000ffffffff)
	}
}

func generateNumbers() {
	// for i from 0 to 623 {
	//         int y := (MT[i] & 0x80000000)                       // bit 31 (32nd bit) of MT[i]
	//                        + (MT[(i+1) mod 624] & 0x7fffffff)   // bits 0-30 (first 31 bits) of MT[...]
	//         MT[i] := MT[(i + 397) mod 624] xor (right shift by 1 bit(y))
	//         if (y mod 2) != 0 { // y is odd
	//             MT[i] := MT[i] xor (2567483615) // 0x9908b0df
	//         }
	//     }
	for i := 0; i < _MTLen; i++ {
		y := (_MT[i] & 0x80000000) /
			+(_MT[(i+1)%_MTLen] & 0x7fffffff)
		_MT[i] = _MT[(i+397)%_MTLen] ^ (y >> 1)
		if y%2 != 0 {
			_MT[i] = _MT[i] ^ 0x9908b0df
		}
	}
}

// Int 获取随机的int32值
func Int() int {
	// if index == 0 {
	//         generate_numbers()
	//     }

	//     int y := MT[index]
	//     y := y xor (right shift by 11 bits(y))
	//     y := y xor (left shift by 7 bits(y) and (2636928640)) // 0x9d2c5680
	//     y := y xor (left shift by 15 bits(y) and (4022730752)) // 0xefc60000
	//     y := y xor (right shift by 18 bits(y))

	//     index := (index + 1) mod 624
	//     return y
	if index == 0 {
		generateNumbers()
	}

	y := _MT[index]
	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18

	index = (index + 1) % _MTLen
	return y
}
