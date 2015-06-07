package main

var BOWSER_8_DOORS = []uint{0x1D4, 0x1D3, 0x1D2, 0x1D1, 0x1CF, 0x1CE, 0x1CD, 0x1CC}

var BOWSER_DARKROOM_ID uint = 0x1BD

func randomizeBowser8Doors(random Random, rom []byte) {
	rooms := []*room{}
	for i := 0; i < len(BOWSER_8_DOORS); i++ {
		id := BOWSER_8_DOORS[i]
		exits := getScreenExits(id, rom)
		rooms = append(rooms, &room{exits[0], id, backupSublevel(id, rom)})
	}
	var j int
	for i := 1; i < len(rooms); i++ {
		j = random.NextInt(i)
		t := rooms[j]
		rooms[j] = rooms[i]
		rooms[i] = t
	}

}

func randomizeBowserEntrances(random Random, _map map[uint]uint, rom []byte, opt *options) {
	backupData(BOWSER_ENTRANCES, rom)
	shuffle(BOWSER_ENTRANCES, random, opt)

	for i := 0; i < len(BOWSER_ENTRANCES); i++ {
		performCopy(BOWSER_ENTRANCES[i], _map, rom, opt)
	}
}

func generateGauntlet(random Random, length uint, rom []byte) {
	if length > uint(len(BOWSER_8_DOORS)) {
		length = uint(len(BOWSER_8_DOORS))
	}

	// get a list of rooms
	rooms := BOWSER_8_DOORS[:uint(len(BOWSER_8_DOORS))-length]

	var j int
	for i := 1; i < len(rooms); i++ {
		j = random.NextInt(i)
		t := rooms[j]
		rooms[j] = rooms[i]
		rooms[i] = t
	}

	numrooms := len(rooms)
	rooms = append(rooms, BOWSER_DARKROOM_ID)

	// copy the first room into both castle entrances
	for _, e := range BOWSER_ENTRANCES {
		copySublevel(e.id, rooms[0], rom)
	}

	// chain together all the rooms \("v")/
	for i := 0; i < numrooms; i++ {
		exits := getScreenExits(rooms[i], rom)
		for j := 0; j < len(exits); j++ {
			rom[exits[j].addr+3] = byte(rooms[i+1] & 0xFF)
		}
	}
}
