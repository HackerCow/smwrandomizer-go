package main

import (
	//"fmt"
	"math"
)

type Random struct {
	Seed int
}

func (r *Random) NextFloat() float64 {
	ret := math.Sin(float64(r.Seed)) * 10000
	r.Seed++
	//fmt.Println("seed in rand=", r.Seed)
	return ret - math.Floor(ret)
}

func (r *Random) NextInt(z int) int {
	return int(r.NextFloat() * float64(z))
}

func (r *Random) NextIntRange(a int, b int) int {
	return a + r.NextInt(b-a)
}
