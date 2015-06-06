package main

import (
	"fmt"
	"math/rand"
)

const (
	NO_CASTLE   = 0
	NORTH_CLEAR = 1
	NORTH_PATH  = 2
)

type stage struct {
	name   string
	world  int
	exits  int
	castle int
	palace int
	ghost  int
	water  int
	id     int
	cpath  int
	tile   [2]int
	out    []string
}

var smwStages = []stage{
	{"yi1", 1, 1, 0, 0, 0, 0, 0x105, NO_CASTLE, [2]int{0x4, 0x28}, []string{"yswitch"}},
	{"yi2", 1, 1, 0, 0, 0, 0, 0x106, NORTH_PATH, [2]int{0xa, 0x28}, []string{"yi3"}},
	{"yi3", 1, 1, 0, 0, 0, 1, 0x103, NORTH_CLEAR, [2]int{0xa, 0x26}, []string{"yi4"}},
	{"yi4", 1, 1, 0, 0, 0, 0, 0x102, NO_CASTLE, [2]int{0xc, 0x24}, []string{"c1"}},
	{"dp1", 2, 2, 0, 0, 0, 0, 0x15, NORTH_PATH, [2]int{0x5, 0x11}, []string{"dp2", "ds1"}},
	{"dp2", 2, 2, 0, 0, 0, 0, 0x9, NORTH_PATH, [2]int{0x3, 0xd}, []string{"dgh", "gswitch"}},
	{"dp3", 2, 1, 0, 0, 0, 0, 0x5, NORTH_CLEAR, [2]int{0x9, 0xa}, []string{"dp4"}},
	{"dp4", 2, 1, 0, 0, 0, 0, 0x6, NO_CASTLE, [2]int{0xb, 0xc}, []string{"c2"}},
	{"ds1", 2, 2, 0, 0, 0, 2, 0xa, NO_CASTLE, [2]int{0x5, 0xe}, []string{"dgh", "dsh"}},
	{"ds2", 2, 1, 0, 0, 0, 0, 0x10b, NORTH_CLEAR, [2]int{0x11, 0x21}, []string{"dp3"}},
	{"vd1", 3, 2, 0, 0, 0, 0, 0x11a, NORTH_CLEAR, [2]int{0x6, 0x32}, []string{"vd2", "vs1"}},
	{"vd2", 3, 2, 0, 0, 0, 1, 0x118, NO_CASTLE, [2]int{0x9, 0x30}, []string{"vgh", "rswitch"}},
	{"vd3", 3, 1, 0, 0, 0, 0, 0x10a, NO_CASTLE, [2]int{0xd, 0x2e}, []string{"vd4"}},
	{"vd4", 3, 1, 0, 0, 0, 0, 0x119, NORTH_PATH, [2]int{0xd, 0x30}, []string{"c3"}},
	{"vs1", 3, 2, 0, 0, 0, 0, 0x109, NO_CASTLE, [2]int{0x4, 0x2e}, []string{"vs2", "sw2"}},
	{"vs2", 3, 1, 0, 0, 0, 0, 0x1, NORTH_CLEAR, [2]int{0xc, 0x3}, []string{"vs3"}},
	{"vs3", 3, 1, 0, 0, 0, 1, 0x2, NORTH_CLEAR, [2]int{0xe, 0x3}, []string{"vfort"}},
	{"cba", 4, 2, 0, 0, 0, 0, 0xf, NORTH_CLEAR, [2]int{0x14, 0x5}, []string{"cookie", "soda"}},
	{"soda", 4, 1, 0, 0, 0, 2, 0x11, NO_CASTLE, [2]int{0x14, 0x8}, []string{"sw3"}},
	{"cookie", 4, 1, 0, 0, 0, 0, 0x10, NORTH_CLEAR, [2]int{0x17, 0x5}, []string{"c4"}},
	{"bb1", 4, 1, 0, 0, 0, 0, 0xc, NO_CASTLE, [2]int{0x14, 0x3}, []string{"bb2"}},
	{"bb2", 4, 1, 0, 0, 0, 0, 0xd, NO_CASTLE, [2]int{0x16, 0x3}, []string{"c4"}},
	{"foi1", 5, 2, 0, 0, 0, 0, 0x11e, NORTH_PATH, [2]int{0x9, 0x37}, []string{"foi2", "fgh"}},
	{"foi2", 5, 2, 0, 0, 0, 1, 0x120, NO_CASTLE, [2]int{0xb, 0x3a}, []string{"foi3", "bswitch"}},
	{"foi3", 5, 2, 0, 0, 0, 0, 0x123, NORTH_CLEAR, [2]int{0x9, 0x3c}, []string{"fgh", "c5"}},
	{"foi4", 5, 2, 0, 0, 0, 0, 0x11f, NORTH_PATH, [2]int{0x5, 0x3a}, []string{"foi2", "fsecret"}},
	{"fsecret", 5, 1, 0, 0, 0, 0, 0x122, NORTH_PATH, [2]int{0x5, 0x3c}, []string{"ffort"}},
	{"ci1", 6, 1, 0, 0, 0, 0, 0x22, NO_CASTLE, [2]int{0x18, 0x16}, []string{"cgh"}},
	{"ci2", 6, 2, 0, 0, 0, 0, 0x24, NORTH_PATH, [2]int{0x15, 0x1b}, []string{"ci3", "csecret"}},
	{"ci3", 6, 2, 0, 0, 0, 0, 0x23, NO_CASTLE, [2]int{0x13, 0x1b}, []string{"ci3", "cfort"}},
	{"ci4", 6, 1, 0, 0, 0, 0, 0x1d, NORTH_PATH, [2]int{0xf, 0x1d}, []string{"ci5"}},
	{"ci5", 6, 1, 0, 0, 0, 0, 0x1c, NORTH_PATH, [2]int{0xc, 0x1d}, []string{"c6"}},
	{"csecret", 6, 1, 0, 0, 0, 0, 0x117, NORTH_CLEAR, [2]int{0x18, 0x29}, []string{"c6"}},
	{"vob1", 7, 1, 0, 0, 0, 0, 0x116, NORTH_CLEAR, [2]int{0x1c, 0x27}, []string{"vob2"}},
	{"vob2", 7, 2, 0, 0, 0, 0, 0x115, NORTH_PATH, [2]int{0x1a, 0x27}, []string{"bgh", "bfort"}},
	{"vob3", 7, 1, 0, 0, 0, 0, 0x113, NORTH_PATH, [2]int{0x15, 0x27}, []string{"vob4"}},
	{"vob4", 7, 2, 0, 0, 0, 0, 0x10f, NORTH_PATH, [2]int{0x15, 0x25}, []string{"sw5"}},
	{"c1", 1, 1, 1, 0, 0, 0, 0x101, NORTH_PATH, [2]int{0xa, 0x22}, []string{"dp1"}},
	{"c2", 2, 1, 2, 0, 0, 0, 0x7, NORTH_PATH, [2]int{0xd, 0xc}, []string{"vd1"}},
	{"c3", 3, 1, 3, 0, 0, 0, 0x11c, NORTH_PATH, [2]int{0xd, 0x32}, []string{"cba"}},
	{"c4", 4, 1, 4, 0, 0, 0, 0xe, NORTH_CLEAR, [2]int{0x1a, 0x3}, []string{"foi1"}},
	{"c5", 5, 1, 5, 0, 0, 0, 0x20, NORTH_CLEAR, [2]int{0x18, 0x12}, []string{"ci1"}},
	{"c6", 6, 1, 6, 0, 0, 0, 0x1a, NORTH_PATH, [2]int{0xc, 0x1b}, []string{"sgs"}},
	{"c7", 7, 1, 7, 0, 0, 0, 0x110, NORTH_PATH, [2]int{0x18, 0x25}, []string{"BOWSER"}},
	{"vfort", 3, 1, -1, 0, 0, 1, 0xb, NORTH_CLEAR, [2]int{0x10, 0x3}, []string{"bb1"}},
	{"ffort", 5, 1, -1, 0, 0, 0, 0x1f, NORTH_CLEAR, [2]int{0x16, 0x10}, []string{"sw4"}},
	{"cfort", 6, 1, -1, 0, 0, 0, 0x1b, NORTH_CLEAR, [2]int{0xf, 0x1b}, []string{"ci4"}},
	{"bfort", 7, 1, -1, 0, 0, 0, 0x111, NORTH_PATH, [2]int{0x1a, 0x25}, []string{"BOWSER"}},
	{"dgh", 2, 2, 0, 0, 1, 0, 0x4, NO_CASTLE, [2]int{0x5, 0xa}, []string{"topsecret", "dp3"}},
	{"dsh", 2, 2, 0, 0, 1, 0, 0x13, NO_CASTLE, [2]int{0x7, 0x10}, []string{"ds2", "sw1"}},
	{"vgh", 3, 1, 0, 0, 1, 0, 0x107, NORTH_CLEAR, [2]int{0x9, 0x2c}, []string{"vd3"}},
	{"fgh", 5, 2, 0, 0, 1, 0, 0x11d, NORTH_CLEAR, [2]int{0x7, 0x37}, []string{"foi1", "foi4"}},
	{"cgh", 6, 1, 0, 0, 1, 0, 0x21, NORTH_CLEAR, [2]int{0x15, 0x16}, []string{"ci2"}},
	{"sgs", 6, 1, 0, 0, 2, 1, 0x18, NORTH_PATH, [2]int{0xe, 0x17}, []string{"vob1"}},
	{"bgh", 7, 2, 0, 0, 1, 0, 0x114, NORTH_PATH, [2]int{0x18, 0x27}, []string{"vob3", "c7"}},
	{"sw1", 8, 2, 0, 0, 0, 0, 0x134, NO_CASTLE, [2]int{0x15, 0x3a}, []string{"sw1", "sw2"}},
	{"sw2", 8, 2, 0, 0, 0, 1, 0x130, NO_CASTLE, [2]int{0x16, 0x38}, []string{"sw2", "sw3"}},
	{"sw3", 8, 2, 0, 0, 0, 0, 0x132, NO_CASTLE, [2]int{0x1a, 0x38}, []string{"sw3", "sw4"}},
	{"sw4", 8, 2, 0, 0, 0, 0, 0x135, NO_CASTLE, [2]int{0x1b, 0x3a}, []string{"sw4", "sw5"}},
	{"sw5", 8, 2, 0, 0, 0, 0, 0x136, NO_CASTLE, [2]int{0x18, 0x3b}, []string{"sw1", "sp1"}},
	{"sp1", 9, 1, 0, 0, 0, 0, 0x12a, NORTH_CLEAR, [2]int{0x14, 0x33}, []string{"sp2"}},
	{"sp2", 9, 1, 0, 0, 0, 0, 0x12b, NORTH_CLEAR, [2]int{0x17, 0x33}, []string{"sp3"}},
	{"sp3", 9, 1, 0, 0, 0, 0, 0x12c, NORTH_CLEAR, [2]int{0x1a, 0x33}, []string{"sp4"}},
	{"sp4", 9, 1, 0, 0, 0, 0, 0x12d, NORTH_CLEAR, [2]int{0x1d, 0x33}, []string{"sp5"}},
	{"sp5", 9, 1, 0, 0, 0, 0, 0x128, NORTH_CLEAR, [2]int{0x1d, 0x31}, []string{"sp6"}},
	{"sp6", 9, 1, 0, 0, 0, 1, 0x127, NORTH_CLEAR, [2]int{0x1a, 0x31}, []string{"sp7"}},
	{"sp7", 9, 1, 0, 0, 0, 0, 0x126, NORTH_CLEAR, [2]int{0x17, 0x31}, []string{"sp8"}},
	{"sp8", 9, 1, 0, 0, 0, 0, 0x125, NORTH_CLEAR, [2]int{0x14, 0x31}, []string{"yi2"}},
	{"yswitch", 1, 0, 0, 1, 0, 0, 0x14, NO_CASTLE, [2]int{0x2, 0x11}, []string{}},
	{"gswitch", 2, 0, 0, 4, 0, 0, 0x8, NO_CASTLE, [2]int{0x1, 0xd}, []string{}},
	{"rswitch", 3, 0, 0, 3, 0, 0, 0x11b, NO_CASTLE, [2]int{0xb, 0x32}, []string{}},
	{"bswitch", 5, 0, 0, 2, 0, 0, 0x121, NO_CASTLE, [2]int{0xd, 0x3a}, []string{}},
	{"topsecret", 2, 0, 0, 0, 0, 0, 0x3, NO_CASTLE, [2]int{0x5, 0x8}, []string{}},
}

func setSlice(dst []byte, offset int, src []byte) {
	for i := offset; i < len(src); i++ {
		dst[i] = src
	}
}

// randomizes slippery/water/tide flags
func randomizeFlags(random Random, stages []stage, buffer []byte) {
	setSlice(buffer, 0x2D8B9, []byte{
		0xA5, 0x0E, 0x8D, 0x0B, 0x01, 0x0A, 0x18, 0x65, 0x0E, 0xA8,
		0xB9, 0x00, 0xE0, 0x85, 0x65, 0xB9, 0x01, 0xE0, 0x85, 0x66,
		0xB9, 0x00, 0xE6, 0x85, 0x68, 0xB9, 0x01, 0xE6, 0x85, 0x69,
		0x80, 0x07})

	setSlice(buffer, 0x026D5, []byte{
		0x20, 0xF5, 0xF9, 0xC9, 0x00, 0xF0, 0x04, 0xC9, 0x05, 0xD0, 0x36})

	setSlice(buffer, 0x02734, []byte{0xEA, 0xEA})

	setSlice(buffer, 0x079F5, []byte{
		0xC2, 0x10, 0xAE, 0x0B, 0x01, 0xBF, 0xE0, 0xFD, 0x03, 0x29,
		0xF0, 0x85, 0x86, 0xBF, 0xE0, 0xFD, 0x03, 0x29, 0x01, 0x85,
		0x85, 0xE2, 0x10, 0xAD, 0x2A, 0x19, 0x60,
	})

	setSlice(buffer, 0x01501, []byte{0x20, 0x91, 0xB0})
	setSlice(buffer, 0x03091, []byte{
		0x64, 0x85, 0x64, 0x86, 0xAD, 0xC6, 0x13, 0x60})
}

func randomizeRom(buffer []byte, seed int64, randomize95Exit bool) {
	stages := make([]stage, len(smwStages))
	copy(stages, smwStages)

	if randomize95Exit {
		// remove dgh and topsecret from rotation
		for i := len(stages) - 1; i >= 0; i-- {
			if stages[i].id == 0x3 || stages[i].id == 0x4 {
				// remove from slice
				stages = append(stages[:i], stages[i+1:]...)
			}
		}
	}

	random := rand.NewSource(seed)
	seedStr := fmt.Sprintf("0x%08x", seed)

	if len(buffer) == 0x80200 {
		buffer = buffer[0x200:]
	}

	randomizeFlags(random, stages, buffer)

}

func main() {
	randomizeRom(nil, 12345, false)
}
