package main

import (
	"errors"
	"fmt"
)

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
		if len(exits) > 0 {
			//fmt.Println("exits", exits[0])
		} else {
			//fmt.Println("exits empty!")
		}
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
	//fmt.Println("ids", ids, "baseid", baseid)
	return ids
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
		//printMD5(rom, "copySublevel")
	}
}

func moveSublevel(to uint, from uint, rom []byte) {
	// copy the sublevel data first
	copySublevel(to, from, rom)
	//fmt.Println("moveSublevel from", to, "from", from)

	// copy the TEST level into the now-freed sublevel slot
	setSlice(rom, LAYER1_OFFSET+3*from, []byte{0x00, 0x80, 0x06})
}

// ASSUMES the layer1 pointer has already been copied to this stage
func fixSublevels(stg *stage, remap map[uint]uint, rom []byte) {
	sublevels := make([]uint, len(stg.copyfrom.sublevels))
	copy(sublevels, stg.copyfrom.sublevels)
	//fmt.Println(sublevels)
	sublevels[0] = stg.id
	//printMD5(rom, "fixSublevels before for 1")
	for i := 1; i < len(sublevels); i++ {
		id := sublevels[i]
		//fmt.Println("id =", id, "stg.id =", stg.id)
		if (id & 0x100) != (stg.id & 0x100) {
			remap[id] = findOpenSublevel(stg.id&0x100, rom)
			newid := remap[id]
			moveSublevel(newid, id, rom)
		}
	}
	//printMD5(rom, "fixSublevels after for 1")

	// fix all screen exits
	var secid uint
	secexitcleanup := []uint{}
	//printMD5(rom, "fixSublevels before for 2")
	for i := 0; i < len(stg.copyfrom.allexits); i++ {
		x := stg.copyfrom.allexits[i]
		target := getSublevelFromExit(x, rom)
		for ele, _ := range remap {
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
	//printMD5(rom, "fixSublevels after for 2")

	//printMD5(rom, "fixSublevels before for 3")
	for _, secid := range secexitcleanup {
		rom[SEC_EXIT_OFFSET_LO+secid] = 0x00
		rom[SEC_EXIT_OFFSET_HI+secid] &= 0xF7
	}
	//printMD5(rom, "fixSublevels after for 3")
	//fmt.Println()
}
