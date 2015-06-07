package main

func removeYoshi(rom []byte, stages []*stage) {
	// change yoshi blocks to 1-up blocks
	for i := 0; i < 36; i++ {
		if rom[0x07080+i] == 0x18 {
			rom[0x07080+i] = 0x0A
		}
	}

	// when baby yoshi grows, he loses all interaction with everything
	rom[0x0A2C1] = 0x02
	rom[0x1C067] = 0x02

	fixDemo(rom)
}

func removeAllPowerups(rom []byte, stages []*stage) {
	powerups := []byte{0x74, 0x75, 0x77}
	removeCape(rom, stages)

	// change powerup blockcodes to multicouns
	blockcodes := []byte{0x02, 0x04, 0x05, 0x08, 0x09}

	for i := 0; i < 36; i++ {
		for _, e := range blockcodes {
			if e == rom[0x07080+i] {
				rom[0x07080+i] = 0x0E
				break
			}
		}
	}

	// change roulette block exclusively to star
	setSlice(rom, 0x0C313, []byte{0x76, 0x76, 0x76, 0x76})

	// change item codes for taking an item through the goal to moving coin
	for i := 0; i < 28; i++ {
		for _, e := range powerups {
			if e == rom[0x07ADF+i] {
				rom[0x07ADF+i] = 0x21
				break
			}
		}
	}

	var id uint
	// remove all fixed powerups in every sublevel
	for id = 0; id < 0x200; id++ {
		sprites := getSpritesBySublevel(id, rom)
		deleteSprites(powerups, sprites, rom)
	}

	// change contents of flying [?]s
	setSlice(rom, 0x0AE88, []byte{0x06, 0x06, 0x06, 0x05, 0x06, 0x06, 0x06, 0x05})

	// remove invisible mushrooms (STZ $14C8,x : RTS)
	setSlice(rom, 0x1C30F, []byte{0x9E, 0xC8, 0x14, 0x60})

	// midpoint shouldn't make you large
	rom[0x072E2] = 0x80

	// yoshi berries (never produce mushroom)
	rom[0x0F0F0] = 0x80

	// remove the powerup peach throws
	rom[0x1A8E9] = 0x00
	rom[0x1A8E4] = 0x00

}

func removeCape(rom []byte, stages []*stage) {
	blockcodes := map[byte]byte{0x08: 0x04, 0x09: 0x05}

	// change feather blockcodes to flower blockcodes
	for i := 0; i < 36; i++ {
		for k, _ := range blockcodes {
			if rom[0x07080+i] == k {
				rom[0x07080+i] = blockcodes[rom[0x07080+i]]
			}
		}
	}

	// change the cape in the roulette block to flower
	rom[0x0C313+2] = 0x75

	// remove capes from super koopas
	setSlice(rom, 0x16AF2, []byte{0xEA, 0xEA, 0xEA, 0xEA})
	setSlice(rom, 0x16B19, []byte{0xEA, 0xEA, 0xEA, 0xEA})

	// change item codes for taking an item through the goal
	for i := 0; i < 28; i++ {
		if rom[0x07ADF+i] == 0x77 {
			rom[0x07ADF+i] = 0x75
		}
	}

	var id uint
	// remove all fixed capes in every sublevel
	for id = 0; id < 0x200; id++ {
		sprites := getSpritesBySublevel(id, rom)
		deleteSprites([]byte{0x77}, sprites, rom)
	}

	setSlice(rom, 0x0AE88, []byte{0x06, 0x02, 0x02, 0x05, 0x06, 0x01, 0x01, 0x05})
}

func randomizePowerups(random Random, rom []byte, stages []*stage) {
	powerups := []byte{0x74, 0x75, 0x77}
	blockmap := map[byte][]byte{
		0x28: {0x28, 0x29},
		0x29: {0x28, 0x29},

		0x30: {0x30, 0x31},
		0x31: {0x30, 0x31},
	}

	var id uint
	for id = 0; id < 0x200; id++ {
		sprites := getSpritesBySublevel(id, rom)
		for i := 0; i < len(sprites); i++ {
			// if we find a bare powerup, replace it with a random powerup
			for _, e := range powerups {
				if e == byte(sprites[i].spriteid) {
					rom[sprites[i].addr+2] = powerups[random.NextInt(len(powerups))]
					break
				}
			}

			// if we find a flying [?], adjust X value (to change its contents)
			if (sprites[i].spriteid == 0x83 || sprites[i].spriteid == 0x84) && (rom[sprites[i].addr+1]&0x30 == 0x10 || rom[sprites[i].addr+1]&0x30 == 0x20) && random.NextInt(2) != 0 {
				rom[sprites[i].addr+1] ^= 0x30
			}

			// FIND AND REPLACE ?/TURN BLOCKS

			start := LAYER1_OFFSET + 3*id
			snes := getPointer(start, 3, rom)

			addr := snesAddressToOffset(snes) + 5
			for ; ; addr += 3 {
				// 0xFF sentinel represents end of level data
				if rom[addr] == 0xFF {
					break
				}

				// pattern looks like the start of the screen exits list
				if (rom[addr]&0xE0) == 0x00 && (rom[addr+1]&0xF5) == 0x00 && rom[addr+2] == 0x00 {
					break
				}

				// pattern looks like a block we want to change
				if (rom[addr]&0x60) == 0x00 && (rom[addr+1]&0xF0) == 0x00 {
					for k, _ := range blockmap {
						if k == rom[addr+2] {
							valid := blockmap[rom[addr+2]]
							rom[addr+2] = valid[random.NextInt(len(valid))]
							break
						}
					}
				}
			}
		}
	}

	// speed up roulette sprite (remove 0xEAs to slow down)
	setSlice(rom, 0x0C327, []byte{0xEA, 0xEA, 0xEA, 0xEA, 0xEA, 0xEA})

	// yoshi can poop any powerup he wants to, don't judge him
	rom[0x0F0F6] = powerups[random.NextInt(len(powerups))]

	// randomize the powerup Peach throws
	rom[0x1A8EE] = powerups[random.NextInt(len(powerups))]
}
