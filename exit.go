package main

import (
	"errors"
	"fmt"
)

type exit struct {
	id     uint
	addr   uint
	screen uint
	water  bool
	issecx bool
	target uint
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

func swapExits(stg *stage, rom []byte) {
	if stg.exits != 2 {
		//fmt.Println("swapExits stg.exits != 2 (", stg.exits, ")")
		return
	}

	// swap output stages
	outa := stg.out[0]
	stg.out[0] = stg.out[1]
	stg.out[1] = outa

	//fmt.Println("stg.out=", stg.out)

	// secret exit triggers event+1
	ndxa := rom[TRANSLEVEL_EVENTS+stg.translevel]
	ndxb := ndxa + 1
	ndxc := ndxa + 2

	// swap the exit directions (nnss----)
	dirs := rom[0x25678+stg.translevel]
	dhi := dirs & 0xC0
	dlo := dirs & 0x30
	rom[0x25678+stg.translevel] = (dirs & 0x0F) | (dhi >> 2) | (dlo << 2)

	//fmt.Println("exit directions=", rom[0x25678+stg.translevel])

	// LAYER 1 ------------------------------

	// "flash" data
	r1flasha := make([]byte, 2)
	copy(r1flasha, rom[uint(0x2585D)+uint(ndxa*2):uint(0x2585D)+uint(ndxa*2)+2])
	r1flashb := make([]byte, 2)
	copy(r1flashb, rom[uint(0x2585D)+uint(ndxb*2):uint(0x2585D)+uint(ndxb*2)+2])

	setSlice(rom, uint(0x2585D)+uint(ndxb*2), r1flasha)
	setSlice(rom, uint(0x2585D)+uint(ndxa*2), r1flashb)

	printMD5(rom, "r1flash")

	// reveal data
	r1reveala := make([]byte, 2)
	copy(r1reveala, rom[uint(0x2593D)+uint(ndxa*2):uint(0x2593D)+uint(ndxa*2)+2])
	r1revealb := make([]byte, 2)
	copy(r1revealb, rom[uint(0x2593D)+uint(ndxb*2):uint(0x2593D)+uint(ndxb*2)+2])

	setSlice(rom, uint(0x2593D)+uint(ndxb*2), r1reveala)
	setSlice(rom, uint(0x2593D)+uint(ndxa*2), r1revealb)

	printMD5(rom, "r1reveal")

	// update offscreen event map
	for i, xor := 0, ndxa^ndxb; i < 44; i++ {
		if rom[0x268E4+i] == ndxa || rom[0x268E4+i] == ndxb {
			rom[0x268E4+i] ^= xor
		}
	}

	printMD5(rom, "offscreen event map")

	// LAYER 2 ------------------------------

	// get offsets into the event data table
	offseta := uint(rom[uint(0x26359)+uint(ndxa)*2]) | (uint(rom[uint(0x26359)+uint(ndxa)*2+1]) << 8)
	//var offsetb = rom[0x26359 + ndxb * 2] | (rom[0x26359 + ndxb * 2 + 1] << 8);
	offsetb := uint(rom[uint(0x26359)+uint(ndxb)*2]) | (uint(rom[uint(0x26359)+uint(ndxb)*2+1]) << 8)
	offsetc := uint(rom[uint(0x26359)+uint(ndxc)*2]) | (uint(rom[uint(0x26359)+uint(ndxc)*2+1]) << 8)

	// get the size of each event
	/*asz*/ _ = offsetb - offseta
	bsz := offsetc - offsetb

	// copy the event data to temporary storage
	eventa := make([]byte, (uint(0x25D8D)+offsetb*4)-(uint(0x25D8D)+offseta*4))
	copy(eventa, rom[uint(0x25D8D)+(offseta*4):uint(0x25D8D)+(offsetb*4)])
	eventb := make([]byte, uint(0x25D8D+(offsetc*4))-(uint(0x25D8D)+offsetb*4))
	copy(eventb, rom[(uint(0x25D8D)+offsetb*4):(uint(0x25D8D)+offsetc*4)])

	//fmt.Println("offsets:", offseta, offsetb, offsetc)
	//fmt.Println("boundaries:", uint(0x25D8D)+offseta*4, uint(0x25D8D)+offsetb*4, uint(0x25D8D)+offsetb*4, uint(0x25D8D)+offsetc*4)
	//fmt.Println("ndx:", ndxa, ndxb, ndxc)

	//fmt.Println(eventa, eventb)
	// update the new offset for where event+1 should go
	offsetb = offseta + bsz

	// copy the event data back into the event table
	setSlice(rom, uint(0x25D8D)+uint(offseta)*4, eventb)
	setSlice(rom, uint(0x25D8D)+uint(offsetb)*4, eventa)

	printMD5(rom, "event data")

	//fmt.Println(offsetb&0xFF, (offsetb>>8)&0xFF)
	// update the offset for event+1 back into the table
	setSlice(rom, uint(0x26359)+uint(ndxb)*2, []byte{byte(offsetb & 0xFF), byte((offsetb >> 8) & 0xFF)})

	printMD5(rom, "offset for event+1")
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

				target := uint(rom[addr+3]) | (id & 0x100)
				x := &exit{id, addr, uint(rom[addr] & 0x1F), (rom[addr+1] & 0x08) != 0, (rom[addr+1] & 0x02) != 0, target}
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
	//fmt.Println("snes", snes)
	return getScreenExitsByAddr(snes, rom, id)
}

func getSecondaryExitTarget(xid uint, rom []byte) uint {
	return uint(rom[SEC_EXIT_OFFSET_LO+xid]) | (xid & 0x100)
}
