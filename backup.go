package main

func backupStage(stg *stage, rom []byte) {
	stg.data = make(map[string][]byte)
	for i := 0; i < len(LEVEL_OFFSETS); i++ {
		off := LEVEL_OFFSETS[i]
		start := off.offset + off.bytes*stg.id

		stg.data[off.name] = make([]byte, off.bytes)
		copy(stg.data[off.name], rom[start:start+off.bytes])

		//printMD5(stg.data[off.name], fmt.Sprintf("data %s", off.name))
		//fmt.Println("backup level=", stg.name, "offset=", off.name, rom[start:start+off.bytes], "(", start, ")")
	}

	stg.translevel = getTranslevel(stg.id)

	for i := 0; i < len(TRANS_OFFSETS); i++ {
		off := TRANS_OFFSETS[i]
		start := off.offset + off.bytes*stg.translevel

		stg.data[off.name] = make([]byte, off.bytes)
		copy(stg.data[off.name], rom[start:start+off.bytes])
		//printMD5(stg.data[off.name], fmt.Sprintf("trans %s", off.name))
	}

	stg.sublevels = getRelatedSublevels(stg.id, rom)

	for _, sublevel := range stg.sublevels {
		screenexits := getScreenExits(sublevel, rom)
		for _, x := range screenexits {
			stg.allexits = append(stg.allexits, x)
		}
	}

	if stg.id == 0x024 {
		// coins - room 2
		xs := getScreenExitsByAddr(0x06E9FB, rom, 0)
		for _, x := range xs {
			stg.allexits = append(stg.allexits, x)
		}
		xs = getScreenExitsByAddr(0x06EAB0, rom, 0)
		for _, x := range xs {
			stg.allexits = append(stg.allexits, x)
		}

		// time - room 3
		xs = getScreenExitsByAddr(0x06EB72, rom, 0)
		for _, x := range xs {
			stg.allexits = append(stg.allexits, x)
		}
		xs = getScreenExitsByAddr(0x06EBBE, rom, 0)
		for _, x := range xs {
			stg.allexits = append(stg.allexits, x)
		}

		// yoshi coins - room 4
		xs = getScreenExitsByAddr(0x06EC7E, rom, 0)
		for _, x := range xs {
			stg.allexits = append(stg.allexits, x)
		}
	}

	stg.data["owtile"] = []byte{rom[getOverworldOffset(stg, 0)]}
	//fmt.Println("data owtile:", stg.data["owtile"])

	//fmt.Println()
}

func backupData(stages []*stage, rom []byte) {
	for i := 0; i < len(stages); i++ {
		backupStage(stages[i], rom)
	}
}
