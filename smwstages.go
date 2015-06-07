package main

type stage struct {
	name       string
	world      int
	exits      int
	castle     int
	palace     int
	ghost      int
	water      int
	id         uint
	cpath      int
	tile       [2]uint
	out        []string
	data       map[string][]byte
	translevel uint
	sublevels  []uint
	allexits   []*exit
	copyfrom   *stage
}

var BOWSER_ENTRANCES = []*stage{
	&stage{"frontdoor", 10, -1, 0, 0, 0, 0, 0x10D, -1, [2]uint{0, 0}, nil, nil, 0, nil, nil, nil},
	&stage{"backdoor", 10, -1, 0, 0, 0, 0, 0x10E, -1, [2]uint{0, 0}, nil, nil, 0, nil, nil, nil},
}
var SMW_STAGES = []stage{
	{"yi1", 1, 1, 0, 0, 0, 0, 0x105, NO_CASTLE, [2]uint{0x4, 0x28}, []string{"yswitch"}, nil, 0, nil, nil, nil},
	{"yi2", 1, 1, 0, 0, 0, 0, 0x106, NORTH_PATH, [2]uint{0xa, 0x28}, []string{"yi3"}, nil, 0, nil, nil, nil},
	{"yi3", 1, 1, 0, 0, 0, 1, 0x103, NORTH_CLEAR, [2]uint{0xa, 0x26}, []string{"yi4"}, nil, 0, nil, nil, nil},
	{"yi4", 1, 1, 0, 0, 0, 0, 0x102, NO_CASTLE, [2]uint{0xc, 0x24}, []string{"c1"}, nil, 0, nil, nil, nil},
	{"dp1", 2, 2, 0, 0, 0, 0, 0x15, NORTH_PATH, [2]uint{0x5, 0x11}, []string{"dp2", "ds1"}, nil, 0, nil, nil, nil},
	{"dp2", 2, 2, 0, 0, 0, 0, 0x9, NORTH_PATH, [2]uint{0x3, 0xd}, []string{"dgh", "gswitch"}, nil, 0, nil, nil, nil},
	{"dp3", 2, 1, 0, 0, 0, 0, 0x5, NORTH_CLEAR, [2]uint{0x9, 0xa}, []string{"dp4"}, nil, 0, nil, nil, nil},
	{"dp4", 2, 1, 0, 0, 0, 0, 0x6, NO_CASTLE, [2]uint{0xb, 0xc}, []string{"c2"}, nil, 0, nil, nil, nil},
	{"ds1", 2, 2, 0, 0, 0, 2, 0xa, NO_CASTLE, [2]uint{0x5, 0xe}, []string{"dgh", "dsh"}, nil, 0, nil, nil, nil},
	{"ds2", 2, 1, 0, 0, 0, 0, 0x10b, NORTH_CLEAR, [2]uint{0x11, 0x21}, []string{"dp3"}, nil, 0, nil, nil, nil},
	{"vd1", 3, 2, 0, 0, 0, 0, 0x11a, NORTH_CLEAR, [2]uint{0x6, 0x32}, []string{"vd2", "vs1"}, nil, 0, nil, nil, nil},
	{"vd2", 3, 2, 0, 0, 0, 1, 0x118, NO_CASTLE, [2]uint{0x9, 0x30}, []string{"vgh", "rswitch"}, nil, 0, nil, nil, nil},
	{"vd3", 3, 1, 0, 0, 0, 0, 0x10a, NO_CASTLE, [2]uint{0xd, 0x2e}, []string{"vd4"}, nil, 0, nil, nil, nil},
	{"vd4", 3, 1, 0, 0, 0, 0, 0x119, NORTH_PATH, [2]uint{0xd, 0x30}, []string{"c3"}, nil, 0, nil, nil, nil},
	{"vs1", 3, 2, 0, 0, 0, 0, 0x109, NO_CASTLE, [2]uint{0x4, 0x2e}, []string{"vs2", "sw2"}, nil, 0, nil, nil, nil},
	{"vs2", 3, 1, 0, 0, 0, 0, 0x1, NORTH_CLEAR, [2]uint{0xc, 0x3}, []string{"vs3"}, nil, 0, nil, nil, nil},
	{"vs3", 3, 1, 0, 0, 0, 1, 0x2, NORTH_CLEAR, [2]uint{0xe, 0x3}, []string{"vfort"}, nil, 0, nil, nil, nil},
	{"cba", 4, 2, 0, 0, 0, 0, 0xf, NORTH_CLEAR, [2]uint{0x14, 0x5}, []string{"cookie", "soda"}, nil, 0, nil, nil, nil},
	{"soda", 4, 1, 0, 0, 0, 2, 0x11, NO_CASTLE, [2]uint{0x14, 0x8}, []string{"sw3"}, nil, 0, nil, nil, nil},
	{"cookie", 4, 1, 0, 0, 0, 0, 0x10, NORTH_CLEAR, [2]uint{0x17, 0x5}, []string{"c4"}, nil, 0, nil, nil, nil},
	{"bb1", 4, 1, 0, 0, 0, 0, 0xc, NO_CASTLE, [2]uint{0x14, 0x3}, []string{"bb2"}, nil, 0, nil, nil, nil},
	{"bb2", 4, 1, 0, 0, 0, 0, 0xd, NO_CASTLE, [2]uint{0x16, 0x3}, []string{"c4"}, nil, 0, nil, nil, nil},
	{"foi1", 5, 2, 0, 0, 0, 0, 0x11e, NORTH_PATH, [2]uint{0x9, 0x37}, []string{"foi2", "fgh"}, nil, 0, nil, nil, nil},
	{"foi2", 5, 2, 0, 0, 0, 1, 0x120, NO_CASTLE, [2]uint{0xb, 0x3a}, []string{"foi3", "bswitch"}, nil, 0, nil, nil, nil},
	{"foi3", 5, 2, 0, 0, 0, 0, 0x123, NORTH_CLEAR, [2]uint{0x9, 0x3c}, []string{"fgh", "c5"}, nil, 0, nil, nil, nil},
	{"foi4", 5, 2, 0, 0, 0, 0, 0x11f, NORTH_PATH, [2]uint{0x5, 0x3a}, []string{"foi2", "fsecret"}, nil, 0, nil, nil, nil},
	{"fsecret", 5, 1, 0, 0, 0, 0, 0x122, NORTH_PATH, [2]uint{0x5, 0x3c}, []string{"ffort"}, nil, 0, nil, nil, nil},
	{"ci1", 6, 1, 0, 0, 0, 0, 0x22, NO_CASTLE, [2]uint{0x18, 0x16}, []string{"cgh"}, nil, 0, nil, nil, nil},
	{"ci2", 6, 2, 0, 0, 0, 0, 0x24, NORTH_PATH, [2]uint{0x15, 0x1b}, []string{"ci3", "csecret"}, nil, 0, nil, nil, nil},
	{"ci3", 6, 2, 0, 0, 0, 0, 0x23, NO_CASTLE, [2]uint{0x13, 0x1b}, []string{"ci3", "cfort"}, nil, 0, nil, nil, nil},
	{"ci4", 6, 1, 0, 0, 0, 0, 0x1d, NORTH_PATH, [2]uint{0xf, 0x1d}, []string{"ci5"}, nil, 0, nil, nil, nil},
	{"ci5", 6, 1, 0, 0, 0, 0, 0x1c, NORTH_PATH, [2]uint{0xc, 0x1d}, []string{"c6"}, nil, 0, nil, nil, nil},
	{"csecret", 6, 1, 0, 0, 0, 0, 0x117, NORTH_CLEAR, [2]uint{0x18, 0x29}, []string{"c6"}, nil, 0, nil, nil, nil},
	{"vob1", 7, 1, 0, 0, 0, 0, 0x116, NORTH_CLEAR, [2]uint{0x1c, 0x27}, []string{"vob2"}, nil, 0, nil, nil, nil},
	{"vob2", 7, 2, 0, 0, 0, 0, 0x115, NORTH_PATH, [2]uint{0x1a, 0x27}, []string{"bgh", "bfort"}, nil, 0, nil, nil, nil},
	{"vob3", 7, 1, 0, 0, 0, 0, 0x113, NORTH_PATH, [2]uint{0x15, 0x27}, []string{"vob4"}, nil, 0, nil, nil, nil},
	{"vob4", 7, 2, 0, 0, 0, 0, 0x10f, NORTH_PATH, [2]uint{0x15, 0x25}, []string{"sw5"}, nil, 0, nil, nil, nil},
	{"c1", 1, 1, 1, 0, 0, 0, 0x101, NORTH_PATH, [2]uint{0xa, 0x22}, []string{"dp1"}, nil, 0, nil, nil, nil},
	{"c2", 2, 1, 2, 0, 0, 0, 0x7, NORTH_PATH, [2]uint{0xd, 0xc}, []string{"vd1"}, nil, 0, nil, nil, nil},
	{"c3", 3, 1, 3, 0, 0, 0, 0x11c, NORTH_PATH, [2]uint{0xd, 0x32}, []string{"cba"}, nil, 0, nil, nil, nil},
	{"c4", 4, 1, 4, 0, 0, 0, 0xe, NORTH_CLEAR, [2]uint{0x1a, 0x3}, []string{"foi1"}, nil, 0, nil, nil, nil},
	{"c5", 5, 1, 5, 0, 0, 0, 0x20, NORTH_CLEAR, [2]uint{0x18, 0x12}, []string{"ci1"}, nil, 0, nil, nil, nil},
	{"c6", 6, 1, 6, 0, 0, 0, 0x1a, NORTH_PATH, [2]uint{0xc, 0x1b}, []string{"sgs"}, nil, 0, nil, nil, nil},
	{"c7", 7, 1, 7, 0, 0, 0, 0x110, NORTH_PATH, [2]uint{0x18, 0x25}, []string{"BOWSER"}, nil, 0, nil, nil, nil},
	{"vfort", 3, 1, -1, 0, 0, 1, 0xb, NORTH_CLEAR, [2]uint{0x10, 0x3}, []string{"bb1"}, nil, 0, nil, nil, nil},
	{"ffort", 5, 1, -1, 0, 0, 0, 0x1f, NORTH_CLEAR, [2]uint{0x16, 0x10}, []string{"sw4"}, nil, 0, nil, nil, nil},
	{"cfort", 6, 1, -1, 0, 0, 0, 0x1b, NORTH_CLEAR, [2]uint{0xf, 0x1b}, []string{"ci4"}, nil, 0, nil, nil, nil},
	{"bfort", 7, 1, -1, 0, 0, 0, 0x111, NORTH_PATH, [2]uint{0x1a, 0x25}, []string{"BOWSER"}, nil, 0, nil, nil, nil},
	{"dgh", 2, 2, 0, 0, 1, 0, 0x4, NO_CASTLE, [2]uint{0x5, 0xa}, []string{"topsecret", "dp3"}, nil, 0, nil, nil, nil},
	{"dsh", 2, 2, 0, 0, 1, 0, 0x13, NO_CASTLE, [2]uint{0x7, 0x10}, []string{"ds2", "sw1"}, nil, 0, nil, nil, nil},
	{"vgh", 3, 1, 0, 0, 1, 0, 0x107, NORTH_CLEAR, [2]uint{0x9, 0x2c}, []string{"vd3"}, nil, 0, nil, nil, nil},
	{"fgh", 5, 2, 0, 0, 1, 0, 0x11d, NORTH_CLEAR, [2]uint{0x7, 0x37}, []string{"foi1", "foi4"}, nil, 0, nil, nil, nil},
	{"cgh", 6, 1, 0, 0, 1, 0, 0x21, NORTH_CLEAR, [2]uint{0x15, 0x16}, []string{"ci2"}, nil, 0, nil, nil, nil},
	{"sgs", 6, 1, 0, 0, 2, 1, 0x18, NORTH_PATH, [2]uint{0xe, 0x17}, []string{"vob1"}, nil, 0, nil, nil, nil},
	{"bgh", 7, 2, 0, 0, 1, 0, 0x114, NORTH_PATH, [2]uint{0x18, 0x27}, []string{"vob3", "c7"}, nil, 0, nil, nil, nil},
	{"sw1", 8, 2, 0, 0, 0, 0, 0x134, NO_CASTLE, [2]uint{0x15, 0x3a}, []string{"sw1", "sw2"}, nil, 0, nil, nil, nil},
	{"sw2", 8, 2, 0, 0, 0, 1, 0x130, NO_CASTLE, [2]uint{0x16, 0x38}, []string{"sw2", "sw3"}, nil, 0, nil, nil, nil},
	{"sw3", 8, 2, 0, 0, 0, 0, 0x132, NO_CASTLE, [2]uint{0x1a, 0x38}, []string{"sw3", "sw4"}, nil, 0, nil, nil, nil},
	{"sw4", 8, 2, 0, 0, 0, 0, 0x135, NO_CASTLE, [2]uint{0x1b, 0x3a}, []string{"sw4", "sw5"}, nil, 0, nil, nil, nil},
	{"sw5", 8, 2, 0, 0, 0, 0, 0x136, NO_CASTLE, [2]uint{0x18, 0x3b}, []string{"sw1", "sp1"}, nil, 0, nil, nil, nil},
	{"sp1", 9, 1, 0, 0, 0, 0, 0x12a, NORTH_CLEAR, [2]uint{0x14, 0x33}, []string{"sp2"}, nil, 0, nil, nil, nil},
	{"sp2", 9, 1, 0, 0, 0, 0, 0x12b, NORTH_CLEAR, [2]uint{0x17, 0x33}, []string{"sp3"}, nil, 0, nil, nil, nil},
	{"sp3", 9, 1, 0, 0, 0, 0, 0x12c, NORTH_CLEAR, [2]uint{0x1a, 0x33}, []string{"sp4"}, nil, 0, nil, nil, nil},
	{"sp4", 9, 1, 0, 0, 0, 0, 0x12d, NORTH_CLEAR, [2]uint{0x1d, 0x33}, []string{"sp5"}, nil, 0, nil, nil, nil},
	{"sp5", 9, 1, 0, 0, 0, 0, 0x128, NORTH_CLEAR, [2]uint{0x1d, 0x31}, []string{"sp6"}, nil, 0, nil, nil, nil},
	{"sp6", 9, 1, 0, 0, 0, 1, 0x127, NORTH_CLEAR, [2]uint{0x1a, 0x31}, []string{"sp7"}, nil, 0, nil, nil, nil},
	{"sp7", 9, 1, 0, 0, 0, 0, 0x126, NORTH_CLEAR, [2]uint{0x17, 0x31}, []string{"sp8"}, nil, 0, nil, nil, nil},
	{"sp8", 9, 1, 0, 0, 0, 0, 0x125, NORTH_CLEAR, [2]uint{0x14, 0x31}, []string{"yi2"}, nil, 0, nil, nil, nil},
	{"yswitch", 1, 0, 0, 1, 0, 0, 0x14, NO_CASTLE, [2]uint{0x2, 0x11}, []string{}, nil, 0, nil, nil, nil},
	{"gswitch", 2, 0, 0, 4, 0, 0, 0x8, NO_CASTLE, [2]uint{0x1, 0xd}, []string{}, nil, 0, nil, nil, nil},
	{"rswitch", 3, 0, 0, 3, 0, 0, 0x11b, NO_CASTLE, [2]uint{0xb, 0x32}, []string{}, nil, 0, nil, nil, nil},
	{"bswitch", 5, 0, 0, 2, 0, 0, 0x121, NO_CASTLE, [2]uint{0xd, 0x3a}, []string{}, nil, 0, nil, nil, nil},
	{"topsecret", 2, 0, 0, 0, 0, 0, 0x3, NO_CASTLE, [2]uint{0x5, 0x8}, []string{}, nil, 0, nil, nil, nil},
}
