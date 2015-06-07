package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"time"
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
	//fmt.Printf("getPointer(ofF=0x%x, len=0x%x, rom)\n", off, len)
	x := uint(0)
	for i := uint(0); i < len; i++ {
		//fmt.Printf("byte = 0x%x\n", rom[off+i])
		x |= uint(rom[off+i]) << uint(i*8)
	}
	//fmt.Printf("x = 0x%x\n", x)
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

func printMD5(arr []byte) {
	if DEBUG_MD5 {
		printMD5Count++
		fmt.Printf("#%d: %x\n", printMD5Count, md5.Sum(arr))
	}
}

// randomizes slippery/water/tide flags
func randomizeFlags(random Random, stages []*stage, rom []byte, opt *options) {
	setSlice(rom, 0x2D8B9, []byte{
		0xA5, 0x0E, 0x8D, 0x0B, 0x01, 0x0A, 0x18, 0x65, 0x0E, 0xA8,
		0xB9, 0x00, 0xE0, 0x85, 0x65, 0xB9, 0x01, 0xE0, 0x85, 0x66,
		0xB9, 0x00, 0xE6, 0x85, 0x68, 0xB9, 0x01, 0xE6, 0x85, 0x69,
		0x80, 0x07})

	setSlice(rom, 0x026D5, []byte{
		0x20, 0xF5, 0xF9, 0xC9, 0x00, 0xF0, 0x04, 0xC9, 0x05, 0xD0, 0x36})

	setSlice(rom, 0x02734, []byte{0xEA, 0xEA})

	setSlice(rom, 0x079F5, []byte{
		0xC2, 0x10, 0xAE, 0x0B, 0x01, 0xBF, 0xE0, 0xFD, 0x03, 0x29,
		0xF0, 0x85, 0x86, 0xBF, 0xE0, 0xFD, 0x03, 0x29, 0x01, 0x85,
		0x85, 0xE2, 0x10, 0xAD, 0x2A, 0x19, 0x60,
	})

	setSlice(rom, 0x01501, []byte{0x20, 0x91, 0xB0})
	setSlice(rom, 0x03091, []byte{
		0x64, 0x85, 0x64, 0x86, 0xAD, 0xC6, 0x13, 0x60})

	var id uint
	for id = 0; id < 0x200; id++ {
		start := uint(LAYER1_OFFSET + 3*id)
		snes := getPointer(start, 3, rom)
		addr := snesAddressToOffset(snes)
		//fmt.Println("addr: ", addr, "; len(rom): ", len(rom))
		numscreens := rom[addr] & 0x1F
		entr := (rom[0x2F200+id] >> 3) & 0x7
		tide := (rom[0x2F200+id] >> 6) & 0x3

		// get default flag setting for the sublevel
		var flag byte
		if entr == 5 {
			flag = 0x80
		} else if entr == 7 {
			flag = 0x01
		} else {
			flag = 0
		}

		// base water on how many screens the stage has
		var toggleWater bool
		if flag&0x01 != 0 {
			// has water
			if opt.delwater {
				toggleWater = true
			}
		} else {
			// doesn't water
			if opt.addwater {
				toggleWater = true
			}
		}
		if random.NextInt(int(math.Max(float64(numscreens), 4))) == 0 && toggleWater {
			flag ^= 0x01
		}

		// force certain stages to not have water
		for _, stage := range NO_WATER_STAGES {
			if id == uint(stage) || tide == 0x1 || tide == 0x2 {
				flag &= 0xF0
				break
			}
		}

		// randomize slippery stages
		if opt.slippery {
			// 12.5%
			if random.NextInt(8) == 0 {
				flag ^= 0x80
			}

			// if the stage is slippery, 33% of the time, changed to "half-slippery"
			if (flag&0x80 != 0) && (random.NextInt(3) == 0) {
				flag ^= 0x90
			}

			if (id == 0xC7) && (flag&0xF0 != 0) {
				fixDemo(rom)
			}
		}
		rom[FLAGBASE+id] = flag
	}
}

func getTranslevel(id uint) uint {
	if id < 0x100 {
		return id
	}
	return id - 0xDC
}

func getScreenExitsByAddr(snes uint, rom []byte, id uint) []*exit {
	var exits []*exit
	addr := snesAddressToOffset(snes) + 5

	for ; ; addr += 3 {
		if rom[addr] == 0xFF {
			break
		}
		if (rom[addr]&0xE0) == 0x00 && (rom[addr+1]&0xF5) == 0x00 && rom[addr+2] == 0x00 {
			for ; ; addr += 4 {
				if rom[addr] == 0xFF {
					break
				}

				x := &exit{id, addr, uint(rom[addr] & 0x1F), (rom[addr+1] & 0x08) != 0, (rom[addr+1] & 0x02) != 0, uint(rom[addr+3]) | (id & 0x100)}
				exits = append(exits, x)
			}
			break
		}
	}
	return exits
}

func getScreenExits(id uint, rom []byte) []*exit {
	start := uint(LAYER1_OFFSET + 3*id)
	snes := getPointer(start, 3, rom)
	return getScreenExitsByAddr(snes, rom, id)
}

func getSecondaryExitTarget(xid uint, rom []byte) uint {
	return uint(rom[SEC_EXIT_OFFSET_LO+xid]) | (xid & 0x100)
}

func getSublevelFromExit(exit *exit, rom []byte) uint {
	if !exit.issecx {
		return exit.target
	}
	return getSecondaryExitTarget(exit.target, rom)
}

func getRelatedSublevels(baseid uint, rom []byte) []uint {
	todo := []uint{baseid}
	ids := []uint{}
	var id uint
	for len(todo) != 0 {
		id = todo[0]
		todo = todo[1:]

		var contains bool
		for _, el := range ids {
			if id == el {
				contains = true
			}
		}
		if contains {
			continue
		}

		ids = append(ids, id)

		exits := getScreenExits(id, rom)
		for i := 0; i < len(exits); i++ {
			x := exits[i]
			next := getSublevelFromExit(x, rom)

			var contains bool
			for _, el := range ids {
				if next == el {
					contains = true
				}
			}
			if !contains {
				todo = append(todo, next)
			}
		}
	}
	return ids
}

func getOverworldOffset(stg *stage, castletop uint) uint {
	tile := stg.tile
	x := tile[0]
	y := tile[1] - castletop

	section := uint((y>>4)*2 + (x >> 4))

	var i uint
	if y >= 0x20 {
		i = 1
	} else {
		i = 0
	}

	return 0x0677DF + section*0x100 + (y&0xF)*0x010 + (x & 0xF) - i
}

func backupStage(stg *stage, rom []byte) {
	stg.data = make(map[string][]byte)
	for i := 0; i < len(LEVEL_OFFSETS); i++ {
		off := LEVEL_OFFSETS[i]
		start := off.offset + off.bytes*stg.id
		stg.data[off.name] = rom[start : start+off.bytes]
		fmt.Println(off.name, "=", rom[start:start+off.bytes], "(", start, ")")
	}

	stg.translevel = getTranslevel(stg.id)

	for i := 0; i < len(TRANS_OFFSETS); i++ {
		off := TRANS_OFFSETS[i]
		start := off.offset + off.bytes*stg.translevel
		stg.data[off.name] = rom[start : start+off.bytes]
	}

	stg.sublevels = getRelatedSublevels(stg.id, rom)

	var allexits []*exit
	for _, sublevel := range stg.sublevels {
		screenexits := getScreenExits(sublevel, rom)
		for _, x := range screenexits {
			allexits = append(allexits, x)
		}
	}

	if stg.id == 0x024 {
		// coins - room 2
		// XXX: maybe 0 isn't right
		xs := getScreenExitsByAddr(0x06E9FB, rom, 0)
		for _, x := range xs {
			allexits = append(allexits, x)
		}
		xs = getScreenExitsByAddr(0x06EAB0, rom, 0)
		for _, x := range xs {
			allexits = append(allexits, x)
		}

		// time - room 3
		xs = getScreenExitsByAddr(0x06EB72, rom, 0)
		for _, x := range xs {
			allexits = append(allexits, x)
		}
		xs = getScreenExitsByAddr(0x06EBBE, rom, 0)
		for _, x := range xs {
			allexits = append(allexits, x)
		}

		// yoshi coins - room 4
		xs = getScreenExitsByAddr(0x06EC7E, rom, 0)
		for _, x := range xs {
			allexits = append(allexits, x)
		}
	}
	stg.allexits = allexits

	// XXX: i don't think this is neccesarry?
	if len(stg.tile) != 0 {
		// XXX: maybe 0 isn't right
		stg.data["owtile"] = []byte{rom[getOverworldOffset(stg, 0)]}
	}
}

func backupData(stages []*stage, rom []byte) {
	for i := 0; i < len(stages); i++ {
		backupStage(stages[i], rom)
	}
}

func shuffle(stages []*stage, random Random, opt *options) {
	rndstages := make([]*stage, len(stages))
	copy(rndstages, stages)

	if opt.randomizeStages {
		var j int
		for i := 1; i < len(rndstages); i++ {
			j = random.NextInt(i)
			t := rndstages[j]
			rndstages[j] = rndstages[i]
			rndstages[i] = t
		}
		for i := 0; i < len(stages); i++ {
			if &rndstages[i] == nil {
				fmt.Println("nil")
			}
			stages[i].copyfrom = rndstages[i]
		}
	}
}

func backupSublevel(id uint, rom []byte) map[string][]byte {
	data := map[string][]byte{}
	for i := 0; i < len(LEVEL_OFFSETS); i++ {
		o := LEVEL_OFFSETS[i]
		x := o.offset + id*o.bytes
		data[o.name] = rom[x : x+o.bytes]
	}
	return data
}

func findOpenSublevel(bank uint, rom []byte) uint {
	bank &= 0x100

	start := []uint{0x025, 0x13C}[bank>>16]
	for i := start; i <= 0xFF; i++ {
		x := bank | i
		os := LAYER1_OFFSET + 3*x
		p := rom[os : os+3]

		// check for TEST level pointer
		if p[0] == 0x00 && p[1] == 0x80 && p[2] == 0x06 {
			return x
		}
	}
	// please god, this should never happen!
	err := fmt.Sprintf("No free sublevels in bank %x", bank)
	panic(errors.New(err))
	return 0
}

func copySublevel(to uint, from uint, rom []byte) {
	// copy all of the level pointers
	for _, o := range LEVEL_OFFSETS {
		fmx := o.offset + from*o.bytes
		tox := o.offset + to*o.bytes
		setSlice(rom, tox, rom[fmx:fmx+o.bytes])
	}
}

func moveSublevel(to uint, from uint, rom []byte) {
	// copy the sublevel data first
	copySublevel(to, from, rom)

	// copy the TEST level into the now-freed sublevel slot
	setSlice(rom, LAYER1_OFFSET+3*from, []byte{0x00, 0x80, 0x06})
}

func findOpenSecondaryExit(bank uint, rom []byte) uint {
	bank &= 0x100
	var i uint
	for i = 0x01; i <= 0xFF; i++ {
		if rom[SEC_EXIT_OFFSET_LO+(bank|i)] == 0x00 {
			return (bank | i)
		}
	}
	// please god, this should never happen!
	err := fmt.Sprintf("No free secondary exits in bank %x", bank)
	panic(errors.New(err))
	return 0
}

// ASSUMES the layer1 pointer has already been copied to this stage
func fixSublevels(stg *stage, remap map[uint]uint, rom []byte) {
	sublevels := make([]uint, len(stg.copyfrom.sublevels))
	copy(sublevels, stg.copyfrom.sublevels)
	sublevels[0] = stg.id

	for i := 1; i < len(sublevels); i++ {
		id := sublevels[i]
		if (id & 0x100) != (stg.id & 0x100) {
			remap[id] = findOpenSublevel(stg.id&0x100, rom)
			newid := remap[id]
			moveSublevel(newid, id, rom)
		}
	}

	// fix all screen exits
	var secid uint
	secexitcleanup := []uint{}
	for i := 0; i < len(stg.copyfrom.allexits); i++ {
		x := stg.copyfrom.allexits[i]
		target := getSublevelFromExit(x, rom)
		for _, ele := range remap {
			if ele == target {
				newtarget := remap[target]
				if !x.issecx {
					rom[x.addr+3] = byte(newtarget)
				} else {
					secid = x.target
					secexitcleanup = append(secexitcleanup, secid)
					newsecid := findOpenSecondaryExit(stg.id&0x100, rom)
					rom[x.addr+3] = byte(newsecid & 0xFF)

					// copy all secondary exit tables
					rom[SEC_EXIT_OFFSET_LO+newsecid] = rom[SEC_EXIT_OFFSET_LO+secid]
					rom[SEC_EXIT_OFFSET_HI+newsecid] = rom[SEC_EXIT_OFFSET_HI+secid]
					rom[SEC_EXIT_OFFSET_X1+newsecid] = rom[SEC_EXIT_OFFSET_X1+secid]
					rom[SEC_EXIT_OFFSET_X2+newsecid] = rom[SEC_EXIT_OFFSET_X2+secid]

					// fix secondary exit target
					rom[SEC_EXIT_OFFSET_LO+newsecid] = byte(newtarget & 0xFF)
					rom[SEC_EXIT_OFFSET_HI+newsecid] &= 0xF7
					rom[SEC_EXIT_OFFSET_HI+newsecid] |= byte((newtarget & 0x100) >> 5)
				}
				break
			}
		}
	}

	for _, secid := range secexitcleanup {
		rom[SEC_EXIT_OFFSET_LO+secid] = 0x00
		rom[SEC_EXIT_OFFSET_HI+secid] &= 0xF7
	}
}

func getRevealedTile(tile uint) uint {
	lookup := map[uint]uint{
		0x7B: 0x7C, 0x7D: 0x7E, 0x76: 0x77, 0x78: 0x79,
		0x7F: 0x80, 0x59: 0x58, 0x57: 0x5E, 0x7A: 0x63,
	}
	var contains bool
	for k, _ := range lookup {
		if k == tile {
			contains = true
		}
	}

	if contains {
		return lookup[tile]
	} else {
		if tile >= 0x66 && tile <= 0x6D {
			return tile + 8
		} else {
			return tile
		}
	}
}

func isCastle(stg *stage) bool {
	return stg.castle > 0
}

func isPermanentTile(stg *stage) bool {
	// some specific tiles MUST be permanent tiles, since the game does not trigger the reveal
	permanentTiles := []string{"sw1", "sw2", "sw3", "sw4", "sw5", "sp1", "yi1", "yi2", "foi1"}

	for _, ele := range permanentTiles {
		if stg.name == ele || isCastle(stg) || isCastle(stg.copyfrom) {
			return true
		}
	}

	revealedTiles := []string{"ci1", "ci2", "ci5", "bb1", "bb2"}
	for _, ele := range revealedTiles {
		if stg.name == ele {
			return false
		}
	}

	if stg.tile[0] >= 0x00 && stg.tile[0] < 0x10 && stg.tile[1] >= 0x35 && stg.tile[1] < 0x40 {
		return false
	}

	return stg.copyfrom.ghost == 1 || stg.copyfrom.castle != 0
}

func getPermanentTile(tile uint) uint {
	lookup := []byte{
		0x7B: 0x7C, 0x7D: 0x7E, 0x76: 0x77, 0x78: 0x79,
		0x7F: 0x80, 0x59: 0x58, 0x57: 0x5E, 0x7A: 0x63,
	}

	var contains bool
	for k, _ := range lookup {
		if uint(k) == tile {
			contains = true
		}
	}

	if contains {
		return uint(lookup[tile])
	} else {
		if tile >= 0x6E && tile <= 0x75 {
			return tile - 8
		} else {
			return tile
		}
	}

}

var hasRun bool

func performCopy(stg *stage, _map map[uint]uint, rom []byte, opt *options) {
	if hasRun {
		return
	}
	hasRun = true
	_map[stg.copyfrom.id] = stg.id
	//printMD5(rom)
	var start uint
	var o levelOffset
	for i := 0; i < len(LEVEL_OFFSETS); i++ {
		fmt.Println(o.name)
		fmt.Println(stg.copyfrom.data[o.name])
		o = LEVEL_OFFSETS[i]
		start = o.offset + o.bytes*stg.id
		setSlice(rom, start, stg.copyfrom.data[o.name])
		//fmt.Println("name", o.name)
		//fmt.Printf("offset: 0x%06x\n", o.offset)
		//fmt.Println(stg.copyfrom.data[o.name])
	}
	//printMD5(rom)

	skiprename := opt.levelNames == LEVEL_NAMES_MATCH_OVERWORLD
	for i := 0; i < len(TRANS_OFFSETS); i++ {
		o := TRANS_OFFSETS[i]
		start := o.offset + o.bytes*stg.translevel
		if skiprename && (o.name == "nameptr") {
			continue
		}
		setSlice(rom, start, stg.copyfrom.data[o.name])
	}
	//printMD5(rom)

	// castle destruction sequence translevels
	if stg.copyfrom.castle > 0 {
		rom[0x049A7+stg.copyfrom.castle-1] = byte(stg.translevel)
	}
	//printMD5(rom)

	// dsh translevel
	if stg.copyfrom.id == 0x013 {
		rom[0x04A0C] = byte(stg.translevel)
	}
	//printMD5(rom)

	// ci2 translevel
	if stg.copyfrom.id == 0x024 {
		rom[0x2DAE5] = byte(stg.translevel)
	}
	//printMD5(rom)

	// moles appear slower in yi2
	if stg.copyfrom.id == 0x106 {
		setSlice(rom, 0x0E2F6, []byte{0xAC, 0xBF, 0x13, 0xEA})
		rom[0x0E2FD] = byte(stg.translevel)
	}
	//printMD5(rom)

	// if the stage we are copying from is default "backwards", we should fix all the
	// associated exits since they have their exits unintuitively reversed
	if stg.copyfrom.id == 0x004 || stg.copyfrom.id == 0x11D {
		swapExits(stg.copyfrom, rom)
		swapExits(stg, rom)
	}
	//printMD5(rom)

	// if we move a stage between 0x100 banks, we need to move sublevels
	// screen exits as well might need to be fixed, even if we don't change banks
	fixSublevels(stg, _map, rom)
	//printMD5(rom)

	// update the overworld tile
	if stg.copyfrom.data != nil && stg.copyfrom.data["owtile"] != nil {
		ow := getRevealedTile(uint(stg.copyfrom.data["owtile"][0]))
		if isPermanentTile(stg) {
			ow = getPermanentTile(ow)
		}
		rom[getOverworldOffset(stg, 0)] = byte(ow)

		// castle 1 is copied into the small version of YI on the main map,
		// which means that whatever stage ends up at c1 needs to be copied there as well
		if stg.name == "c1" {
			rom[0x67A54] = byte(ow)
			if !isCastle(stg.copyfrom) {
				rom[0x67A44] = 0x00
			}
		}

		// moving a castle here, need to add a castle top
		if isCastle(stg.copyfrom) && !isCastle(stg) {
			rom[getOverworldOffset(stg, 1)] = []byte{0x00, 0xA6, 0x4C}[stg.cpath]
		}

		// moving a castle away, need to fix the top tile
		if !isCastle(stg.copyfrom) && isCastle(stg) {
			rom[getOverworldOffset(stg, 1)] = []byte{0x00, 0x00, 0x10}[stg.cpath]
		}

		// fix offscreen event tiles
		for k, _ := range OFFSCREEN_EVENT_TILES {
			if k == stg.name {
				x := OFFSCREEN_EVENT_TILES[stg.name]
				tile := getPermanentTile(ow)

				if (stg.name == "c1") && (stg.copyfrom.castle != 0) {
					tile = 0x81
				}
				setSlice(rom, uint(0x26994)+uint(x), []byte{byte(tile), 0x00})
				break
			}
		}
	}
}

func removeAutoscrollers(rom []byte) {
	var id uint
	for id = 0; id < 0x200; id++ {
		// don't remove autoscroller for DP2 due to issues with the
		// layer 2 movement. sorry, but this is a wontfix issue.
		if id == 0x009 {
			continue
		}

		snes := getPointer(uint(LAYER1_OFFSET+3*id), 3, rom)
		addr := snesAddressToOffset(snes) + 4

		var sprites = getSpritesBySublevel(id, rom)

		// fix the v-scroll if we find autoscroller sprites
		for i := 0; i < len(sprites); i++ {
			switch sprites[i].spriteid {
			case 0xE8:
				rom[addr] &= 0xCF
				rom[addr] |= 0x10
				break

			case 0xF3:
				rom[addr] &= 0xCF
				rom[addr] |= 0x00
				break
			default:
				continue
			}
		}

		// remove the actual autoscroller sprites
		deleteSprites([]byte{0xE8, 0xF3}, sprites, rom)
	}
}

func fixOverworldEvents(stages []*stage, rom []byte) {
	_map := map[byte]*stage{}
	var stg *stage
	for i := 0; i < len(stages); i++ {
		if stages[i].copyfrom != nil {
			stg = stages[i].copyfrom
		} else {
			stg = stages[i]
		}

		if stg.copyfrom.palace != 0 || stg.copyfrom.castle != 0 {
			_map[rom[TRANSLEVEL_EVENTS+stg.copyfrom.translevel]] = stg
		}
	}

	EVENTS := 0x265D6
	COORDS := 0x265B6
	VRAM := 0x26587

	for i := 0; i < 16; i++ {
		_ /*stage*/ = _map[rom[EVENTS+i]]
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
		rom[COORDS+i*2+1] = byte((y>>4)*2 + (x >> 4))
		s := rom[COORDS+i*2+1]

		if s >= 0x4 {
			if x == 0 {
				y--
			}
			x = (x + 0xF) & 0xF
		}
		pos := ((y & 0xF) << 4) | (x & 0xF)
		rom[COORDS+i*2] = byte(pos)

		// update vram values for overworld updates
		rom[VRAM+i*2] = byte(0x20 | ((y & 0x10) >> 1) | ((x & 0x10) >> 2) | ((y & 0x0C) >> 2))
		rom[VRAM+i*2+1] = byte(((y & 0x03) << 6) | ((x & 0x0F) << 1))

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

		ghostspritebase := 0x27625 + 35
		ghostsubmapbase := 0x27666

		for i := 0; i < len(stages); i++ {
			if stages[i].copyfrom.ghost == 1 {
				s := stages[i].tile[1] >= 0x20
				x := stages[i].tile[0]*0x10 - 0x10
				if s {
					x -= 0x0010
				}

				y := stages[i].tile[1]*0x10 - 0x08
				if s {
					x -= 0x0200
				}

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
			}
		}
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

func randomizeKoopaKids(_map map[uint]uint, random Random, rom []byte) {
	bossrooms := []bossroom{}
	for i := 0; i < len(KOOPA_KIDS); i++ {
		// find the actual sublevel holding this boss fight
		oldbr := KOOPA_KIDS[i]
		newbr := oldbr
		for k, _ := range _map {
			if k == oldbr {
				newbr = _map[oldbr]
				break
			}
		}

		// save this information
		bossrooms = append(bossrooms, bossroom{uint(i + 1), newbr, 0})
	}

	var j int
	for i := 1; i < len(bossrooms); i++ {
		j = random.NextInt(i)
		t := bossrooms[j]
		bossrooms[j] = bossrooms[i]
		bossrooms[i] = t
	}

	for i := 0; i < len(bossrooms); i++ {
		bossrooms[i].cto = uint(i + 1)
	}

	hold0 := findOpenSublevel(0x000, rom)
	moveSublevel(hold0, bossrooms[0].sublevel, rom)

	for i := 1; i < len(bossrooms); i++ {
		moveSublevel(bossrooms[len(bossrooms)-1].sublevel, hold0, rom)
	}
	moveSublevel(bossrooms[len(bossrooms)-1].sublevel, hold0, rom)

	// TODO: fix castle names
}

func randomizeBossDifficulty(random Random, rom []byte) {
	// health of ludwig+morton+roy
	var jhp = random.NextIntRange(2, 9)
	rom[0x0CFCD] = byte(jhp)
	rom[0x0D3FF] = byte(jhp + 9)

	// health of big boo
	rom[0x181A2] = byte(random.NextIntRange(2, 7))

	// health of wendy+lemmy
	var whp = random.NextIntRange(2, 6)
	rom[0x1CE1A] = byte(whp)
	rom[0x1CED4] = byte(whp - 1)

	// health of bowser phase1, and phases2+3
	var bhp = random.NextIntRange(1, 6)
	rom[0x1A10B] = byte(bhp)
	rom[0x1A683] = byte(bhp)

	// distance iggy+larry slides when hit (jump, fireball)
	rom[0x0FD00] = byte(random.NextIntRange(0x08, 0x30))
	rom[0x0FD46] = byte(random.NextIntRange(0x08, 0x28))

	// randomize reznor
	if 0 == random.NextInt(3) {
		setSlice(rom, 0x198C7, []byte{0x38, 0xE9})
	}
	rom[0x198C9] = byte(random.NextIntRange(0x01, 0x04))
}

func getChecksum(rom []byte) uint16 {
	var checksum uint16 = 0
	for i := 0; i < len(rom); i++ {
		checksum += uint16(rom[i])
		checksum &= 0xFFFF
	}
	return checksum
}

func centerPad(str string, length uint) string {
	for uint(len(str)) < length {
		if len(str)&1 != 0 {
			str = " " + str
		} else {
			str = str + " "
		}
	}
	return str
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

func randomizeRom(rom []byte, seed int, opt *options) {
	rawStages := make([]stage, len(SMW_STAGES))
	copy(rawStages, SMW_STAGES)

	stages := make([]*stage, len(rawStages))
	for i := 0; i < len(rawStages); i++ {
		stages[i] = &rawStages[i]
	}

	if opt.randomize95Exit {
		// remove dgh and topsecret from rotation
		for i := len(stages) - 1; i >= 0; i-- {
			if stages[i].id == 0x3 || stages[i].id == 0x4 {
				// remove from slice
				stages = append(stages[:i], stages[i+1:]...)
			}
		}
	}

	//printMD5(rom)

	random := Random{seed}
	vseed := fmt.Sprintf("0x%08x", seed)

	if len(rom) == 0x80200 {
		rom = rom[0x200:]
	}

	//printMD5(rom)

	// randomize all of the slippery/water flags
	randomizeFlags(random, stages, rom, opt)
	//printMD5(rom)

	// NOTE: MAKE SURE ANY TABLES BACKED UP BY THIS ROUTINE ARE GENERATED *BEFORE*
	// THIS POINT. OTHERWISE, WE ARE BACKING UP UNINITIALIZED MEMORY!
	backupData(stages, rom)
	//printMD5(rom)

	// put all the stages into buckets (any stage in same bucket can be swapped)
	buckets := makeBuckets(stages, opt)
	//printMD5(rom)

	// decide which stages will be swapped with which others
	for i := 0; i < len(buckets); i++ {
		shuffle(buckets[i], random, opt)
	}
	//printMD5(rom)

	// quick stage lookup table
	stagelookup := map[string]*stage{}
	for i := 0; i < len(stages); i++ {
		stagelookup[stages[i].copyfrom.name] = stages[i]
	}

	if opt.randomizeBowserDoors {
		randomizeBowser8Doors(random, rom)
	}
	//printMD5(rom)

	globalremapping := map[uint]uint{}

	switch opt.bowser {
	case BOWSER_SWAP_DOORS:
		randomizeBowserEntrances(random, globalremapping, rom, opt)
		break
	case BOWSER_GAUNTLET:
		generateGauntlet(random, 8, rom)
		break
	case BOWSER_MINI_GAUNTLET:
		generateGauntlet(random, uint(random.NextIntRange(3, 6)), rom)
		break
	default:
		break
	}
	//printMD5(rom)

	switch opt.powerups {
	case POWERUP_RANDOMIZE:
		randomizePowerups(random, rom, stages)
		break
	case POWERUP_NO_CAPE:
		removeCape(rom, stages)
		break
	case POWERUP_SMALL_ONLY:
		removeAllPowerups(rom, stages)
		break
	default:
		break
	}
	//printMD5(rom)

	if opt.noYoshi {
		removeYoshi(rom, stages)
	}
	//printMD5(rom)

	// remove all autoscroller sprites and update v-scroll, if checked
	if opt.removeAutoscrollers {
		removeAutoscrollers(rom)
	}
	//printMD5(rom)

	// update level names if randomized
	if opt.levelNamesCustom {
		randomizeLevelNames(random, rom)
	}
	//printMD5(rom)

	// swap all the level name pointers RIGHT before we perform the copy
	if opt.levelNames == LEVEL_NAMES_RANDOM_STAGE {
		shuffleLevelNames(stages, random)
	}
	//printMD5(rom)

	for i := 0; i < len(stages); i++ {
		performCopy(stages[i], globalremapping, rom, opt)
		// randomly swap the normal/secret exits
		if opt.randomizeExits && random.NextFloat() > 0.5 {
			swapExits(stages[i], rom)
		}
	}
	printMD5(rom)

	// fix castle/fort/switch overworld tile events
	fixOverworldEvents(stages, rom)
	//printMD5(rom)

	// fix Roy/Larry castle block paths
	fixBlockPaths(stagelookup, rom)
	//printMD5(rom)

	// fix message box messages
	fixMessageBoxes(stages, rom)
	//printMD5(rom)

	if opt.randomizeKoopaKids {
		randomizeKoopaKids(globalremapping, random, rom)
	}
	//printMD5(rom)

	if opt.randomizeBossDiff {
		randomizeBossDifficulty(random, rom)
	}
	//printMD5(rom)

	// disable the forced no-yoshi intro on moved stages
	rom[0x2DA1D] = 0x60

	// infinite lives?
	if opt.infiniteLives {
		setSlice(rom, 0x050D8, []byte{0xEA, 0xEA, 0xEA})
	}

	// write version number and the randomizer seed to the rom
	var checksum = fmt.Sprintf("%04x", getChecksum(rom))
	writeToTitle(VERSION_STRING+" @"+vseed+"-"+checksum, 0x2, rom)
	//printMD5(rom)

	fixChecksum(rom)
	printMD5(rom)
}

func checkMD5(data []byte) error {
	currentSum := fmt.Sprintf("%x", md5.Sum(data))
	for _, sum := range ORIGINAL_MD5 {
		if sum == currentSum {
			return nil
		}
	}
	return errors.New("The MD5 sum of the ROM did not match that of an original ROM")
}

func check(err error) {
	if err != nil {
		panic(err)
	}

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	opt, err := parseFlags()
	check(err)

	dat, err := ioutil.ReadFile(opt.filename)
	check(err)

	err = checkMD5(dat)
	check(err)

	var seed int

	if opt.customSeed == "" {
		seed = int(rand.Int31())
	} else {
		seed64, err := strconv.ParseInt(opt.customSeed, 16, 64)
		seed = int(seed64)
		check(err)
	}

	fmt.Printf("Using seed %x\n", seed)

	randomizeRom(dat, seed, opt)
}
