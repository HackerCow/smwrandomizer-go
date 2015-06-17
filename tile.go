package main

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
	lookup := map[byte]byte{
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
		return uint(lookup[byte(tile)])
	} else {
		if tile >= 0x6E && tile <= 0x75 {
			return tile - 8
		} else {
			return tile
		}
	}
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

func getRevealedTile(tile uint) uint {
	lookup := map[byte]byte{
		0x7C: 0x7B, 0x7E: 0x7D, 0x77: 0x76, 0x79: 0x78,
		0x80: 0x7F, 0x58: 0x59, 0x5E: 0x57, 0x63: 0x7A,
	}
	var contains bool
	for k, _ := range lookup {
		if uint(k) == tile {
			contains = true
		}
	}

	if contains {
		//fmt.Println("tile in get", tile)
		return uint(lookup[byte(tile)])
	} else {
		if tile >= 0x66 && tile <= 0x6D {
			return tile + 8
		} else {
			return tile
		}
	}
}
