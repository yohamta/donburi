package helper

import "math/rand"

func RangeFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func Chance(percent float64) bool {
	return rand.Float64() > percent
}
