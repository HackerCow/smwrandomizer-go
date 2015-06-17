package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

// randomizes slippery/water/tide flags
func randomizeFlags(random *Random, stages []*stage, rom []byte, opt *options) {
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

func shuffleLevels(stages []*stage, random *Random, opt *options) []*stage {
	rndstages := make([]*stage, len(stages))
	copy(rndstages, stages)
	if opt.randomizeStages {
		rndstages = shuffle(random, rndstages)
	}

	for i := 0; i < len(stages); i++ {
		stages[i].copyfrom = rndstages[i]
		//fmt.Println(stages[i].copyfrom.sublevels)
		//fmt.Println(rndstages[i].sublevels)
	}

	return stages
}

func performCopy(stg *stage, _map map[uint]uint, rom []byte, opt *options) {
	_map[stg.copyfrom.id] = stg.id
	//printMD5(rom)
	var start uint
	var o levelOffset
	//printMD5(rom, "performCopy 1")
	for i := 0; i < len(LEVEL_OFFSETS); i++ {
		//fmt.Println(o.name)
		//fmt.Println(stg.copyfrom.data[o.name])
		o = LEVEL_OFFSETS[i]
		start = o.offset + o.bytes*stg.id
		//fmt.Println("stage=", stg.name, "copyfrom=", stg.copyfrom.name, "name=", o.name, "start=", start, stg.copyfrom.data[o.name])
		setSlice(rom, start, stg.copyfrom.data[o.name])
		//fmt.Println("name", o.name)
		//fmt.Printf("offset: 0x%06x\n", o.offset)
		//fmt.Println(stg.copyfrom.data[o.name])
	}
	//printMD5(rom)

	//printMD5(rom, "performCopy 2")
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

	//printMD5(rom, "performCopy 3")
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

	//printMD5(rom, "performCopy 4")
	// if the stage we are copying from is default "backwards", we should fix all the
	// associated exits since they have their exits unintuitively reversed
	if stg.copyfrom.id == 0x004 || stg.copyfrom.id == 0x11D {
		//fmt.Println("SWAPPING EXITS D:")
		swapExits(stg.copyfrom, rom)
		swapExits(stg, rom)
	}
	//printMD5(rom)

	// if we move a stage between 0x100 banks, we need to move sublevels
	// screen exits as well might need to be fixed, even if we don't change banks
	fixSublevels(stg, _map, rom)

	//printMD5(rom, "after fixSublevels")

	// update the overworld tile
	if stg.copyfrom.data != nil && stg.copyfrom.data["owtile"] != nil {
		//fmt.Println("tile", uint(stg.copyfrom.data["owtile"][0]))
		ow := getRevealedTile(uint(stg.copyfrom.data["owtile"][0]))
		if isPermanentTile(stg) {
			//fmt.Println("permanent")
			ow = getPermanentTile(ow)
		}
		overworldOffset := getOverworldOffset(stg, 0)
		rom[overworldOffset] = byte(ow)
		//fmt.Println(overworldOffset, ow)
		//printMD5(rom, "getOverworldOffset set")

		// castle 1 is copied into the small version of YI on the main map,
		// which means that whatever stage ends up at c1 needs to be copied there as well
		if stg.name == "c1" {
			rom[0x67A54] = byte(ow)
			if !isCastle(stg.copyfrom) {
				rom[0x67A44] = 0x00
			}
		}
		//printMD5(rom, "c1 fixed")

		// moving a castle here, need to add a castle top
		if isCastle(stg.copyfrom) && !isCastle(stg) {
			rom[getOverworldOffset(stg, 1)] = []byte{0x00, 0xA6, 0x4C}[stg.cpath]
		}
		//printMD5(rom, "castle top")

		// moving a castle away, need to fix the top tile
		if !isCastle(stg.copyfrom) && isCastle(stg) {
			rom[getOverworldOffset(stg, 1)] = []byte{0x00, 0x00, 0x10}[stg.cpath]
		}
		//printMD5(rom, "fix castle top")

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
		//printMD5(rom, "offscreen event tiles")
	} else {
		//fmt.Println("null")
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

	printMD5(rom, "start")

	random := &Random{seed}

	vseed := fmt.Sprintf("%08x", seed)

	if len(rom) == 0x80200 {
		rom = rom[0x200:]
	}

	printMD5(rom, "created")

	// randomize all of the slippery/water flags
	randomizeFlags(random, stages, rom, opt)
	printSeed(random)
	printMD5(rom, "randomizeFlags")

	// NOTE: MAKE SURE ANY TABLES BACKED UP BY THIS ROUTINE ARE GENERATED *BEFORE*
	// THIS POINT. OTHERWISE, WE ARE BACKING UP UNINITIALIZED MEMORY!
	backupData(stages, rom)
	printSeed(random)
	printMD5(rom, "backupData")

	// put all the stages into buckets (any stage in same bucket can be swapped)
	buckets := makeBuckets(stages, opt)
	printSeed(random)
	printMD5(rom, "makeBuckets")

	// decide which stages will be swapped with which others
	for i := 0; i < len(buckets); i++ {
		buckets[i] = shuffleLevels(buckets[i], random, opt)
	}
	printSeed(random)
	printMD5(rom, "shuffle")

	// quick stage lookup table
	stagelookup := map[string]*stage{}
	for i := 0; i < len(stages); i++ {
		stagelookup[stages[i].copyfrom.name] = stages[i]
	}

	if opt.randomizeBowserDoors {
		randomizeBowser8Doors(random, rom)
	}
	printSeed(random)
	printMD5(rom, "bowserDoors")

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
	printSeed(random)
	printMD5(rom, "gauntlet")

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
	printSeed(random)
	printMD5(rom, "powerups")

	if opt.noYoshi {
		removeYoshi(rom, stages)
	}
	printSeed(random)
	printMD5(rom, "removeYoshi")

	// remove all autoscroller sprites and update v-scroll, if checked
	if opt.removeAutoscrollers {
		removeAutoscrollers(rom)
	}
	printSeed(random)
	printMD5(rom, "autoscrollers")

	// update level names if randomized
	if opt.levelNamesCustom {
		randomizeLevelNames(random, rom)
	}
	printSeed(random)
	printMD5(rom, "randomizeLevelNames")

	// swap all the level name pointers RIGHT before we perform the copy
	if opt.levelNames == LEVEL_NAMES_RANDOM_STAGE {
		shuffleLevelNames(stages, random)
	}
	printSeed(random)
	printMD5(rom, "shuffleLevelNames")

	for i := 0; i < len(stages); i++ {
		//fmt.Println("\n\n\n")
		//printMD5(rom, "before performcopy")
		performCopy(stages[i], globalremapping, rom, opt)
		//printMD5(rom, "after performcopy")
		// randomly swap the normal/secret exits
		if opt.randomizeExits && random.NextFloat() > 0.5 {
			//fmt.Println("swapexits")
			swapExits(stages[i], rom)
		}
	}
	printSeed(random)
	printMD5(rom, "performCopy")

	// fix castle/fort/switch overworld tile events
	fixOverworldEvents(stages, rom)
	printSeed(random)
	printMD5(rom, "fixOverworldEvents")

	// fix Roy/Larry castle block paths
	fixBlockPaths(stagelookup, rom)
	printSeed(random)
	printMD5(rom, "fixBlockPaths")

	// fix message box messages
	fixMessageBoxes(stages, rom)
	printSeed(random)
	printMD5(rom, "fixMessageBoxes")

	if opt.randomizeKoopaKids {
		randomizeKoopaKids(globalremapping, random, rom)
	}
	printSeed(random)
	printMD5(rom, "randomizeKoopaKids")

	if opt.randomizeBossDiff {
		randomizeBossDifficulty(random, rom)
	}
	printSeed(random)
	printMD5(rom, "randomizeBossDifficulty")

	// disable the forced no-yoshi intro on moved stages
	rom[0x2DA1D] = 0x60

	// infinite lives?
	if opt.infiniteLives {
		setSlice(rom, 0x050D8, []byte{0xEA, 0xEA, 0xEA})
	}

	// write version number and the randomizer seed to the rom
	var checksum = fmt.Sprintf("%04x", getChecksum(rom))
	writeToTitle(VERSION_STRING+" @"+vseed+"-"+checksum, 0x2, rom)
	printMD5(rom, "writeToTitle")

	fixChecksum(rom)
	printSeed(random)
	printMD5(rom, "fixChecksum")

	var filename string
	if opt.outFilename == "" {
		filename = fmt.Sprintf("smw-%s.sfc", vseed)
	} else {
		filename = opt.outFilename
	}
	ioutil.WriteFile(filename, rom, 0644)
	fmt.Println("Saved ROM as", filename)
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
