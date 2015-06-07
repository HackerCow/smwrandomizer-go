package main

import (
	"math"
)

func isSameBucket(a *stage, b *stage, opt *options) bool {
	// a different number of exits, different bucket
	if a.exits != b.exits {
		return false
	}

	// can only put castles where there is room for them
	if (b.cpath != NO_CASTLE) != (a.cpath != NO_CASTLE) {
		return false
	}

	// if same-type, most both be castle/non-castle for same bucket
	if opt.randomizeSameType && (math.Signbit(float64(a.castle)) != math.Signbit(float64(b.castle))) {
		return false
	}

	// if same-type, most both be ghost/non-ghost for same bucket
	if opt.randomizeSameType && (a.ghost != b.ghost) {
		return false
	}

	// if same-type, most both be water/non-water for same bucket
	if opt.randomizeSameType && (a.water != b.water) {
		return false
	}

	// if same-type, most both be palace/non-palace for same bucket
	if opt.randomizeSameType && math.Signbit(float64(a.palace)) != math.Signbit(float64(b.palace)) {
		return false
	}

	// option: randomize only within worlds
	if opt.randomizeSameWorld && a.world != b.world {
		return false
	}
	return true
}

func makeBuckets(stages []*stage, opt *options) [][]*stage {
	buckets := make([][]*stage, 0)
	for x := 0; x < len(stages); x++ {
		var i int
		var st *stage
		st = stages[x]
		for i = 0; i < len(buckets); i++ {
			if isSameBucket(st, buckets[i][0], opt) {
				buckets[i] = append(buckets[i], st)
				break
			}
		}
		// if we didn't put a stage in a bucket, new bucket
		if i == len(buckets) {
			buckets = append(buckets, []*stage{st})
		}
	}
	return buckets
}
