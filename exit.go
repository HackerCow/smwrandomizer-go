package main

type exit struct {
	id     uint
	addr   uint
	screen uint
	water  bool
	issecx bool
	target uint
}

func swapExits(stg *stage, rom []byte) {
	// swap output stages
	outa := stg.out[0]
	stg.out[0] = stg.out[1]
	stg.out[1] = outa

	// secret exit triggers event+1
	ndxa := rom[TRANSLEVEL_EVENTS+stg.translevel]
	ndxb := ndxa + 1
	ndxc := ndxa + 2

	// swap the exit directions (nnss----)
	dirs := rom[0x25678+stg.translevel]
	dhi := dirs & 0xC0
	dlo := dirs & 0x30
	rom[0x25678+stg.translevel] = (dirs & 0x0F) | (dhi >> 2) | (dlo << 2)

	// LAYER 1 ------------------------------

	// "flash" data
	r1flasha := rom[uint(0x2585D)+uint(ndxa*2) : uint(0x2585D)+uint(ndxa*2+2)]
	r1flashb := rom[uint(0x2585D)+uint(ndxb*2) : uint(0x2585D)+uint(ndxb*2+2)]

	setSlice(rom, uint(0x2585D)+uint(ndxb*2), r1flasha)
	setSlice(rom, uint(0x2585D)+uint(ndxa*2), r1flashb)

	// reveal data
	r1reveala := rom[uint(0x2593D)+uint(ndxa*2) : uint(0x2593D)+uint(ndxa*2+2)]
	r1revealb := rom[uint(0x2593D)+uint(ndxb*2) : uint(0x2593D)+uint(ndxb*2+2)]

	setSlice(rom, uint(0x2593D)+uint(ndxb*2), r1reveala)
	setSlice(rom, uint(0x2593D)+uint(ndxa*2), r1revealb)

	// update offscreen event map
	for i, xor := 0, ndxa^ndxb; i < 44; i++ {
		if rom[0x268E4+i] == ndxa || rom[0x268E4+i] == ndxb {
			rom[0x268E4+i] ^= xor
		}
	}

	// LAYER 2 ------------------------------

	// get offsets into the event data table
	offseta := rom[uint(0x26359)+uint(ndxa*2)] | (rom[uint(0x26359)+uint(ndxa*2+1)] << 8)
	offsetb := rom[uint(0x26359)+uint(ndxb*2)] | (rom[uint(0x26359)+uint(ndxb*2+1)] << 8)
	offsetc := rom[uint(0x26359)+uint(ndxc*2)] | (rom[uint(0x26359)+uint(ndxc*2+1)] << 8)

	// get the size of each event
	/*asz*/ _ = offsetb - offseta
	bsz := offsetc - offsetb

	// copy the event data to temporary storage
	eventa := rom[uint(0x25D8D)+uint(offseta*4) : uint(0x25D8D)+uint(offsetb*4)]
	eventb := rom[uint(0x25D8D)+uint(offsetb*4) : uint(0x25D8D)+uint(offsetc*4)]

	// update the new offset for where event+1 should go
	offsetb = offseta + bsz

	// copy the event data back into the event table
	setSlice(rom, uint(0x25D8D)+uint(offseta*4), eventb)
	setSlice(rom, uint(0x25D8D)+uint(offsetb*4), eventa)

	// update the offset for event+1 back into the table
	setSlice(rom, uint(0x26359)+uint(ndxb*2), []byte{offsetb & 0xFF, (offsetb >> 8) & 0xFF})
}
