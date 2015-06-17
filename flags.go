package main

import (
	"errors"
	"flag"
)

type levelNamesOptions int
type bowserOptions int
type powerupOptions int

const (
	LEVEL_NAMES_MATCH_STAGE levelNamesOptions = 1 + iota
	LEVEL_NAMES_MATCH_OVERWORLD
	LEVEL_NAMES_RANDOM_STAGE

	BOWSER_DEFAULT bowserOptions = 1 + iota
	BOWSER_SWAP_DOORS
	BOWSER_GAUNTLET
	BOWSER_MINI_GAUNTLET

	POWERUP_DEFAULT powerupOptions = 1 + iota
	POWERUP_RANDOMIZE
	POWERUP_NO_CAPE
	POWERUP_SMALL_ONLY
)

type options struct {
	filename    string
	outFilename string
	customSeed  string

	randomizeStages    bool
	randomizeSameWorld bool
	randomizeSameType  bool
	randomize95Exit    bool

	levelNames       levelNamesOptions
	levelNamesCustom bool

	bowser               bowserOptions
	randomizeBowserDoors bool

	powerups powerupOptions
	noYoshi  bool

	slippery bool
	addwater bool
	delwater bool

	randomizeExits      bool
	randomizeKoopaKids  bool
	randomizeBossDiff   bool
	removeAutoscrollers bool

	infiniteLives bool
}

func parseFlags() (*options, error) {
	seed := flag.String("seed", "", "A custom seed to use")
	outFilename := flag.String("out", "", "The filename the ROM should be saved to. Default is \"smw-[seed].sfc\"")
	randomizeStages := flag.Bool("randomizeStages", true, "Randomly Swap Stages (Sorta the point, ya know?)")
	randomizeSameWorld := flag.Bool("randomizeSameWorld", false, "Keep Same World (Randomize only within world)")
	randomizeSameType := flag.Bool("randomizeSameType", false, "Keep Same Type (Ghost, Water, Normal)")
	randomize95Exit := flag.Bool("randomize95Exit", false, "95 Exit Mode (Don't move Donut Ghost House)")

	levelNames := flag.String("levelNames", "matchStage", "Can be one of \"matchStage\", \"matchOverworld\", \"randomStage\"")
	customRandomNames := flag.Bool("customRandomNames", false, "Custom Random Names")

	bowsersCastle := flag.String("bowsersCastle", "default", "Can be one of \"default\", \"swapDoors\", \"gauntlet\", \"miniGauntlet\"")
	randomizeBowserDoors := flag.Bool("randomizeBowserDoors", false, "Randomize Bowser's 1-8 Doors")

	powerups := flag.String("powerups", "default", "Can be one of \"default\", \"randomize\", \"noCape\", \"smallOnly\"")
	noYoshi := flag.Bool("noYoshi", false, "No Yoshi")

	slippery := flag.Bool("slippery", false, "Randomize Slippery/Ice Physics")
	addWater := flag.Bool("addWater", false, "Randomly Add Water to Stages")
	delWater := flag.Bool("delWater", false, "Randonly Remove Water from Stages (Can Be Difficult!)")

	randomizeExits := flag.Bool("randomizeExits", false, "Randomize Exits (Swap Normal/Secret Randomly)")
	randomizeKoopaKids := flag.Bool("randomizeKoopaKids", false, "Shuffle Koopa Kids")
	randomizeBossDiff := flag.Bool("randomizeBossDiff", false, "Randomize Boss Difficulty (Koopa Kids, Big Boo, Reznor, Bowser)")
	removeAutoscrollers := flag.Bool("removeAutoscrollers", false, "Remove Autoscrollers (except DP2)")
	infiniteLives := flag.Bool("infiniteLives", false, "Infinite Lifes")

	debugMD5 := flag.Bool("debugmd5", false, "Print MD5 Debug output")

	flag.Parse()
	args := flag.Args()

	DEBUG_MD5 = *debugMD5

	if len(args) != 1 {
		return nil, errors.New("No/Invalid filename supplied. For usage use -help.")
	}

	filename := args[0]

	var levelNamesOpt levelNamesOptions

	switch *levelNames {
	case "matchStage":
		levelNamesOpt = LEVEL_NAMES_MATCH_STAGE
	case "matchOverworld":
		levelNamesOpt = LEVEL_NAMES_MATCH_OVERWORLD
	case "randomStage":
		levelNamesOpt = LEVEL_NAMES_RANDOM_STAGE
	default:
		return nil, errors.New("Invalid \"levelNames\" option.")
	}

	var bowserOpt bowserOptions

	switch *bowsersCastle {
	case "default":
		bowserOpt = BOWSER_DEFAULT
	case "swapDoors":
		bowserOpt = BOWSER_SWAP_DOORS
	case "gauntlet":
		bowserOpt = BOWSER_GAUNTLET
	case "miniGauntlet":
		bowserOpt = BOWSER_MINI_GAUNTLET
	default:
		return nil, errors.New("Invalid \"bowsersCastle\" option.")
	}

	var powerupOpt powerupOptions

	switch *powerups {
	case "default":
		powerupOpt = POWERUP_DEFAULT
	case "randomize":
		powerupOpt = POWERUP_RANDOMIZE
	case "noCape":
		powerupOpt = POWERUP_NO_CAPE
	case "smallOnly":
		powerupOpt = POWERUP_SMALL_ONLY
	default:
		return nil, errors.New("Invalid \"powerups\" option.")
	}

	opt := options{filename, *outFilename, *seed, *randomizeStages, *randomizeSameWorld, *randomizeSameType,
		*randomize95Exit, levelNamesOpt, *customRandomNames, bowserOpt, *randomizeBowserDoors,
		powerupOpt, *noYoshi, *slippery, *addWater, *delWater, *randomizeExits,
		*randomizeKoopaKids, *randomizeBossDiff, *removeAutoscrollers, *infiniteLives}

	return &opt, nil
}
