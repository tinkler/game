package util

import (
	"math/rand"
	"time"
)

type KrandKind int

const (
	KC_RAND_KIND_NUM KrandKind = iota
	KC_RAND_KIND_LOWER
	KC_RAND_KIND_UPPER
	KC_RAND_KIND_ALL
)

func Krand(size int, kind KrandKind) []byte {
	ikind := int(kind)
	kinds := [][]int{{10, 48}, {26, 97}, {26, 65}}
	result := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if kind > 2 || kind < 0 {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
