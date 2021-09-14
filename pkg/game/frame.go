package game

import "time"

type Frame struct {
	// 同步操作计时器
	S int
	// 同步时间
	T time.Time
	// 同步值
	V interface{}
}
