package main

type sprite struct {
	stage    uint
	addr     uint
	spriteid byte
}

func getSpritesBySublevel(id uint, rom []byte) []*sprite {
	start := SPRITE_OFFSET + 2*id
	snes := 0x070000 | getPointer(start, 2, rom)
	addr := snesAddressToOffset(snes) + 1
	sprites := []*sprite{}

	for ; ; addr += 3 {
		if rom[addr] == 0xFF {
			break
		}
		s := &sprite{id, addr, rom[addr+2]}

		sprites = append(sprites, s)
	}
	return sprites
}

func deleteSprites(todelete []byte, sprites []*sprite, rom []byte) {
	if len(sprites) == 0 {
		return
	}

	length := len(sprites)
	base := sprites[0].addr

	for i := length - 1; i >= 0; i-- {
		for j := 0; j < len(todelete); j++ {
			if todelete[j] == sprites[i].spriteid {
				for j := i + 1; j < length; j++ {
					addr := base + uint(j*3)
					setSlice(rom, addr-3, rom[:addr])
					sprites[j].addr -= 3 // not needed, but correct
				}

				// remove the sprite object from the list
				sprites = append(sprites[:i], sprites[i+1:]...)
				length--
			}
		}
	}
	// end of list marker
	rom[base+uint(length*3)] = 0xFF
}
