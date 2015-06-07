package main

import (
	"math"
)

type Random struct {
	seed int
}

func (r *Random) NextFloat() float64 {
	ret := math.Sin(float64(r.seed)) * 10000
	r.seed++
	return ret - math.Floor(ret)
}

func (r *Random) NextInt(z int) int {
	return int(r.NextFloat() * float64(z))
}

func (r *Random) NextIntRange(a int, b int) int {
	return a + r.NextInt(b-a)
}
