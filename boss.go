package main

func randomizeKoopaKids(_map map[uint]uint, random *Random, rom []byte) {
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

func randomizeBossDifficulty(random *Random, rom []byte) {
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
