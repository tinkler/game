package attr

import "math"

type Offset struct {
	Dx, Dy float64
}

func FromDirection(direction, distance float64) Offset {
	return Offset{distance * math.Cos(direction), distance * math.Sin(direction)}
}

func (current Offset) Add(offset Offset) Offset {
	return Offset{current.Dx + offset.Dx, current.Dy + offset.Dy}
}

func (current Offset) Direction() float64 {
	return math.Atan2(current.Dx, current.Dy)
}

func (current Offset) Distance() float64 {
	return math.Sqrt(current.Dx*current.Dx + current.Dy*current.Dy)
}
