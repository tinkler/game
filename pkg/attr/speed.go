package attr

type Speed float64

func (s Speed) IsZero() bool {
	return float64(s) == 0.0
}
