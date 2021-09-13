package game

import "time"

type Frame struct {
	// step
	S int
	T time.Time
	V interface{}
}
