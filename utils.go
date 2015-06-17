package main

import (
	"crypto/md5"
	"fmt"
	"math"
)

const (
	NO_CASTLE   = 0
	NORTH_CLEAR = 1
	NORTH_PATH  = 2

	FLAGBASE = 0x1FDE0

	SEC_EXIT_OFFSET_LO = 0x2F800
	SEC_EXIT_OFFSET_HI = 0x2FE00
	SEC_EXIT_OFFSET_X1 = 0x2FA00
	SEC_EXIT_OFFSET_X2 = 0x2FC00

	TRANSLEVEL_EVENTS = 0x2D608
)

var LAYER1_OFFSET = uint(0x2E000)
var LAYER2_OFFSET = uint(0x2E600)
var SPRITE_OFFSET = uint(0x2EC00)

type levelOffset struct {
	name   string
	bytes  uint
	offset uint
}

type transOffset struct {
	name   string
	bytes  uint
	offset uint
}

var KOOPA_KIDS = []uint{0x1F6, 0x0E5, 0x1F2, 0x0D9, 0x0CC, 0x0D3, 0x1EB}

type bossroom struct {
	cfm      uint
	sublevel uint
	cto      uint
}

var VERSION_STRING = "v1.2"

var TRANS_OFFSETS = []transOffset{
	// name pointer
	{"nameptr", 2, 0x220FC},
}

var LEVEL_OFFSETS = []levelOffset{
	// layer data
	{"layer1", 3, LAYER1_OFFSET},
	{"layer2", 3, LAYER2_OFFSET},
	{"sprite", 2, SPRITE_OFFSET},

	// secondary header data
	{"header1", 1, 0x2F000},
	{"header2", 1, 0x2F200},
	{"header3", 1, 0x2F400},
	{"header4", 1, 0x2F600},

	// custom data
	{"lvflags", 1, FLAGBASE},
}

var OFFSCREEN_EVENT_TILES = map[string]byte{
	"dp1":     0x00,
	"c1":      0x0A,
	"dp3":     0x10,
	"vs2":     0x20,
	"ds2":     0x26,
	"cba":     0x28,
	"csecret": 0x3A,
	"vob1":    0x3E,
	"yswitch": 0x40,
	"vd1":     0x4C,
}

var DEBUG_MD5 bool

var ORIGINAL_MD5 = []string{"cdd3c8c37322978ca8669b34bc89c804", "dbe1f3c8f3a0b2db52b7d59417891117"}

var NO_WATER_STAGES = []int{0x01A, 0x0DC, 0x111, 0x1CF, 0x134, 0x1F8, 0x0C7, 0x1E3, 0x1E2, 0x1F2, 0x0CC}

func setSlice(dst []byte, offset uint, src []byte) {
	var i uint
	for i = 0; i < uint(len(src)); i++ {
		dst[i+offset] = src[i]
	}
}

func getPointer(off uint, len uint, rom []byte) uint {
	//fmt.Printf("getPointer(off=%d, len=%d, rom)\n", off, len)
	x := uint(0)
	for i := uint(0); i < len; i++ {
		//fmt.Printf("byte = %d\n", rom[off+i])
		x |= uint(rom[off+i]) << uint(i*8)
	}
	//fmt.Printf("x = %d\n", x)
	return x
}

func snesAddressToOffset(addr uint) uint {
	//fmt.Printf("snesAddressToOffset(addr=0x%x)\n", addr)
	ret := ((addr&0xFF0000)>>16)*0x8000 + (addr & 0x00FFFF) - 0x8000
	//fmt.Printf("ret=0x%x\n", ret)
	return ret
}

// end the demo after 10 inputs. this is a safe place to stop both for
// the no-yoshi intro and for slippery/water intros
func fixDemo(rom []byte) {
	for i := 10; i < 34; i++ {
		rom[0x01C1F+i*2] = 0x00
	}
	rom[0x01C1F+34] = 0xFF
}

var printMD5Count = 0

func printMD5(arr []byte, msg string) {
	if DEBUG_MD5 {
		printMD5Count++
		fmt.Printf("#%d: %x (%s)\n", printMD5Count, md5.Sum(arr), msg)
	}
}

func shuffle(r *Random, s []*stage) []*stage {
	for i := 1; i < len(s); i++ {
		j := r.NextInt(i)
		t := s[j]
		s[j] = s[i]
		s[i] = t
	}
	return s
}

func fixOverworldEvents(stages []*stage, rom []byte) {
	_map := map[byte]*stage{}
	for i := 0; i < len(stages); i++ {
		var stg *stage
		if stages[i].copyfrom != nil {
			stg = stages[i].copyfrom
		} else {
			stg = stages[i]
		}

		if stg.copyfrom.palace != 0 || stg.copyfrom.castle != 0 {
			_map[rom[TRANSLEVEL_EVENTS+stg.copyfrom.translevel]] = stg
			//fmt.Println(TRANSLEVEL_EVENTS+stg.copyfrom.translevel, stg)
		}
	}

	EVENTS := 0x265D6
	COORDS := 0x265B6
	VRAM := 0x26587

	for i := 0; i < 16; i++ {
		var stg *stage
		stg = _map[rom[EVENTS+i]]
		if stg == nil || stg.copyfrom == nil {
			continue
		}
		rom[EVENTS+i] = rom[TRANSLEVEL_EVENTS+stg.translevel]
		tile := stg.tile
		x := tile[0]
		y := tile[1]
		if stg.copyfrom.castle > 0 {
			y--
		}

		//fmt.Println("x", x, "y", y)

		rom[COORDS+i*2+1] = byte((y>>4)*2 + (x >> 4))
		s := rom[COORDS+i*2+1]

		if s >= 0x4 {
			if x == 0 {
				y--
			}
			x = (x + 0xF) & 0xF
		}

		//fmt.Println("s", s)
		pos := ((y & 0xF) << 4) | (x & 0xF)
		rom[COORDS+i*2] = byte(pos)

		// update vram values for overworld updates
		rom[VRAM+i*2] = byte(0x20 | ((y & 0x10) >> 1) | ((x & 0x10) >> 2) | ((y & 0x0C) >> 2))
		rom[VRAM+i*2+1] = byte(((y & 0x03) << 6) | ((x & 0x0F) << 1))
	}
	//printMD5(rom, "vram values")

	// always play bgm after a stage
	rom[0x20E2E] = 0x80

	// fix castle crush for castles crossing section numbers
	setSlice(rom, 0x266C5, []byte{0x20, 0x36, 0xA2, 0xEA})
	setSlice(rom, 0x26F22, []byte{0x20, 0x35, 0xA2, 0xEA})

	setSlice(rom, 0x22235, []byte{
		0xA8, 0x29, 0xFF, 0x00, 0xC9, 0xF0, 0x00, 0x98, 0x90,
		0x03, 0x69, 0xFF, 0x00, 0x69, 0x10, 0x00, 0x60})

	setSlice(rom, 0x27625, []byte{0x00, 0x00, 0x01, 0xE0, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00,
		0x06, 0x70, 0x01, 0x20, 0x00,
		0x07, 0x38, 0x00, 0x8A, 0x01,
		0x00, 0x58, 0x00, 0x7A, 0x00,
		0x08, 0x88, 0x01, 0x18, 0x00,
		0x09, 0x48, 0x01, 0xFC, 0xFF})

	//printMD5(rom, "castle crush")

	ghostspritebase := 0x27625 + 35
	ghostsubmapbase := 0x27666

	for i := 0; i < len(stages); i++ {
		if stages[i].copyfrom.ghost == 1 {
			s := stages[i].tile[1] >= 0x20

			sub := 0
			if s {
				sub = 0x0010
			}
			x := int(stages[i].tile[0])*0x10 - 0x10 - sub

			sub = 0
			if s {
				sub = 0x0200
			}

			y := int(stages[i].tile[1])*0x10 - 0x08 - sub

			//fmt.Println("s", s, "x", x, "y", y, stages[i].tile)

			setSlice(rom, uint(ghostspritebase), []byte{
				0x0A, byte(x & 0xFF), byte((x >> 8) & 0xFF),
				byte(y & 0xFF), byte((y >> 8) & 0xFF)})

			ghostspritebase += 5

			if s {
				rom[ghostsubmapbase] = 1
			} else {
				rom[ghostsubmapbase] = 0
			}
			ghostsubmapbase++
			//printMD5(rom, "ghost stuff loop")
		}
	}

	//printMD5(rom, "ghost stuff")
	setSlice(rom, 0x2766F, []byte{0x01, 0x20, 0x40, 0x60, 0x80, 0xA0})
	setSlice(rom, 0x27D7F, []byte{0xF0, 0x02, 0xA9, 0x01, 0x5D, 0x5C, 0xF6, 0xD0, 0x2C, 0x9B, 0x8A, 0x85, 0x0A, 0x0A, 0x0A,
		0x18, 0x65, 0x0A, 0xAA, 0xC2, 0x20, 0xBD, 0x17, 0xF6, 0x18, 0x69, 0x10, 0x00, 0x20, 0xB9,
		0xFF, 0xE2, 0x30, 0xC9, 0x7A, 0xF0, 0x10, 0xBB, 0xA9, 0x34, 0xBC, 0x95, 0x0E, 0x30, 0x02,
		0xA9, 0x44, 0xEB, 0xA9, 0x60, 0x20, 0x06, 0xFB})

	setSlice(rom, 0x27FB1, []byte{0x22, 0x75, 0xF6, 0x04, 0x5C, 0x00, 0x80, 0x7F, 0x85, 0x0A, 0xE2, 0x20, 0x4A, 0x4A, 0x4A,
		0x4A, 0x85, 0x0A, 0xC2, 0x20, 0xBD, 0x19, 0xF6, 0x18, 0x69, 0x08, 0x00, 0x85, 0x0C, 0xE2,
		0x20, 0x29, 0xF0, 0x18, 0x65, 0x0A, 0xEB, 0xA5, 0x0D, 0x0A, 0x65, 0x0B, 0xEB, 0xC2, 0x30,
		0xAA, 0xBF, 0x00, 0xC8, 0x7E, 0x85, 0x58, 0x60})

	rom[0x276B0] = 0x0A
	rom[0x27802] = 0x0A

	setSlice(rom, 0x01AA4, []byte{0x22, 0xB1, 0xFF, 0x04})
}

func fixBlockPaths(lookup map[string]*stage, rom []byte) {
	c5 := lookup["c5"]
	c7 := lookup["c7"]
	hitrans := byte(math.Max(float64(c5.translevel), float64(c7.translevel)))

	// swap some values if roy and larry end up in the wrong order
	if c7.translevel < c5.translevel {
		setSlice(rom, 0x19307, []byte{0xEF, 0x93})
		setSlice(rom, 0x1930C, []byte{0xA4, 0x93})
	}

	/*
		org $0392F8
			LDA $13BF
			NOP #3
			CMP #$xx
	*/
	setSlice(rom, 0x192F8, []byte{0xAD, 0xBF, 0x13, 0xEA, 0xEA, 0xEA, 0xC9, hitrans})
}

func fixMessageBoxes(stages []*stage, rom []byte) {
	// mapping for where translevels moved
	var transmap = map[uint]uint{}
	for i := 0; i < len(stages); i++ {
		transmap[stages[i].copyfrom.translevel] = stages[i].translevel
	}

	// 23 bytes in table at 0x2A590
	for i := 0; i < 23; i++ {
		val := rom[0x2A590+i]
		t := val & 0x7F
		for k, _ := range transmap {
			if uint(t) == k {
				rom[0x2A590+i] = byte(transmap[uint(t)] | uint(val&0x80))
			}
		}
	}
}

func getChecksum(rom []byte) uint16 {
	var checksum uint16 = 0
	for i := 0; i < len(rom); i++ {
		checksum += uint16(rom[i])
		checksum &= 0xFFFF
	}
	return checksum
}

func fixChecksum(rom []byte) {
	var checksum = getChecksum(rom)

	// checksum
	rom[0x7FDE] = byte(checksum & 0xFF)
	rom[0x7FDF] = byte((checksum >> 8) & 0xFF)

	// checksum ^ 0xFFFF
	rom[0x7FDC] = (rom[0x7FDE] & 0xFF) ^ 0xFF
	rom[0x7FDD] = (rom[0x7FDF] & 0xFF) ^ 0xFF
}

func printSeed(r *Random) {
	//fmt.Println("seed=", r.Seed)
}

func getTranslevel(id uint) uint {
	if id < 0x100 {
		return id
	}
	return id - 0xDC
}
