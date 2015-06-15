package main

import (
	"fmt"
	"strings"
)

var TEXT_MAPPING = map[uint16]byte{
	'A': 0x00, 'B': 0x01, 'C': 0x02, 'D': 0x03, 'E': 0x04, 'F': 0x05, 'G': 0x06, 'H': 0x07, 'I': 0x08,
	'J': 0x09, 'K': 0x0a, 'L': 0x0b, 'M': 0x0c, 'N': 0x0d, 'O': 0x0e, 'P': 0x0f, 'Q': 0x10, 'R': 0x11,
	'S': 0x12, 'T': 0x13, 'U': 0x14, 'V': 0x15, 'W': 0x16, 'X': 0x17, 'Y': 0x18, 'Z': 0x19, '!': 0x1a,
	'.': 0x1b, '-': 0x1c, ',': 0x1d, '?': 0x1e, ' ': 0x1f, 'a': 0x40, 'b': 0x41, 'c': 0x42, 'd': 0x43,
	'e': 0x44, 'f': 0x45, 'g': 0x46, 'h': 0x47, 'i': 0x48, 'j': 0x49, 'k': 0x4a, 'l': 0x4b, 'm': 0x4c,
	'n': 0x4d, 'o': 0x4e, 'p': 0x4f, 'q': 0x50, 'r': 0x51, 's': 0x52, 't': 0x53, 'u': 0x54, 'v': 0x55,
	'w': 0x56, 'x': 0x57, 'y': 0x58, 'z': 0x59, '#': 0x5a, '(': 0x5b, ')': 0x5c, '\'': 0x5d, '·': 0x5e,
	'1': 0x64, '2': 0x65, '3': 0x66, '4': 0x67, '5': 0x68, '6': 0x69, '7': 0x6a, '0': 0x6b,

	'\uE032': 0x32, '\uE033': 0x33, '\uE034': 0x34, '\uE035': 0x35, '\uE036': 0x36, '\uE037': 0x37,

	'\uE038': 0x38, '\uE039': 0x39, '\uE03A': 0x3a, '\uE03B': 0x3b, '\uE03C': 0x3c,
}

var TITLE_STRINGS = [][]string{
	[]string{"YOSHI'S ", "MARIO'S ", "LUIGI'S ", "DEATHLY ", "LEMMY'S ", "LARRY'S ", "WENDY'S ", "KOOPA'S "},
	[]string{"STAR ", "HYPE ", "MOON "},
	[]string{"#1 IGGY'S "},
	[]string{"#2 MORTON'S "},
	[]string{"#3 LEMMY'S "},
	[]string{"#4 LUDWIG'S "},
	[]string{"#5 ROY'S "},
	[]string{"#6 WENDY'S "},
	[]string{"#7 LARRY'S "},
	[]string{"DONUT ", "PIZZA ", "DEATH ", "KOOPA ", "FUDGE ", "PLUTO ", "KAIZO ", "SKULL ", "MARIO ", "SUSHI ", "BAGEL ", "BREAD "},
	[]string{"GREEN "},
	[]string{"TOP SECRET AREA ", "TAKE A BREAK ", "WHY THE RUSH? ", "LEVEL OF SHAME ", "KEEP YOUR CAPES "},
	[]string{"VANILLA ", "DIAMOND ", "CALZONE ", "EMERALD ", "BUTTERY ", "DOLPHIN "},
	[]string{"\uE038\uE039\uE03A\uE03B\uE03C "}, // YELLOW
	[]string{"RED "},
	[]string{"BLUE "},
	[]string{"BUTTER BRIDGE ", "CHEESE BRIDGE ", "APPLE ISTHMUS ", "ASIAGO CHEESE "},
	[]string{"CHEESE BRIDGE ", "BUTTER BRIDGE ", "PASTA PLATEAU ", "BOUNCING SAWS "},
	[]string{"SODA LAKE ", "POP OCEAN ", "INK SWAMP "},
	[]string{"COOKIE MOUNTAIN ", "GREEN HILL ZONE ", "WALUIGI LAND ", "PRINCESS VALLEY ", "DINO-RHINO LAND "},
	[]string{"FOREST ", "CANOPY ", "JUNGLE "},
	[]string{"CHOCOLATE ", "CHEEZCAKE ", "PEPPERONI "},
	[]string{"CHOCO-GHOST HOUSE ", "HAUNTED MANSION ", "HOUSE OF HORROR ", "HOUSE OF TERROR "},
	[]string{"SUNKEN GHOST SHIP ", "GHOSTS OF YOSHI ", "SMB3 AIRSHIP "},
	[]string{"VALLEY ", "SUMMIT ", "RIVERS ", "THREAT ", "WOUNDS ", "GALAXY "},
	[]string{"BACK DOOR ", "NO ENTRY ", "GO AWAY ", "LEAVE NOW ", "PISS OFF! ", "BEWARE! "},
	[]string{"FRONT DOOR ", "GET READY ", "FINAL BOSS "},
	[]string{"GNARLY ", "WACKY ", "CRAZY ", "KOOKY ", "NUTTY "},
	[]string{"TUBULAR ", "FUCKYOU ", "(-.-) ", "GETREKT "},
	[]string{"WAY COOL ", "GLORIOUS ", "STYLISH ", "SUAVE "},
	[]string{"HOUSE ", "ABODE ", "CONDO ", "TOWER "},
	[]string{"ISLAND ", "MIRAGE ", "TUNNEL ", "CAVERN ", "BRIDGE ", "GALAXY "},
	[]string{"SWITCH PALACE "},
	[]string{"CASTLE ", "TEMPLE ", "DOMAIN "},
	[]string{"PLAINS ", "TUNDRA ", "MEADOW ", "CAVERN ", "BRIDGE "},
	[]string{"GHOST HOUSE ", "BOO'S HAUNT ", "GRAVEYARD "},
	[]string{"SECRET ", "TEMPLE "},
	[]string{"DOME ", "ZONE ", "HELL "},
	[]string{"FORTRESS ", "DUNGEON ", "DUNGEONS ", "PANTHEON ", "CAPITAL ", "CENTRE ", "CENTER "},
	[]string{"OF\uE032\uE033\uE034\uE035\uE036\uE037ON ", "OF DISDAIN ", "OF VISIONS ", "HAPPY TIME "}, // OF ILLUSION
	[]string{"OF BOWSER ", "OF KOOPAS ", "OF SORROW ", "OF CLOUDS ", "OF SPIKES "},
	[]string{"ROAD ", "WARP ", "PATH ", "ZONE ", "LINE "},
	[]string{"WORLD "},
	[]string{"AWESOME ", "SPOOKY ", "STRANGE ", "AMAZING ", "MYSTERY ", "ENGLAND ", "TWITCHY ", "RADICAL "},
	[]string{"1", "0"},
	[]string{"2", "?"},
	[]string{"3", "X"},
	[]string{"4", "6"},
	[]string{"5", "7"},
	[]string{"PALACE", "TEMPLE", "SHRINE"},
	[]string{"AREA", "ZONE", "SPOT", "HILL"},
	[]string{"GROOVY", "CRAZY", "DEATH!", "CANADA"},
	[]string{"MONDO", "KINKY", "Kappa", "GG!", "????", "OMG!"},
	[]string{"OUTRAGEOUS", "UNNATURAL", "MENTAL", "MADNESS", "TRY NOCAPE", "BibleThump", "FABULOUS"},
	[]string{"FUNKY", "GREAT", "WEIRD", "BINGO"},
	[]string{"HOUSE", "ABODE", "CONDO", "TOWER"},
	[]string{" "},
}

func shuffleLevelNames(stages []*stage, random *Random) {
	ptrs := make([]byte, len(stages))
	for i, e := range stages {
		ptrs[i] = e.data["nameptr"][0]
	}

	var j int
	for i := 1; i < len(ptrs); i++ {
		j = random.NextInt(i)
		t := ptrs[j]
		ptrs[j] = ptrs[i]
		ptrs[i] = t
	}

	for i := 0; i < len(ptrs); i++ {
		stages[i].data["nameptr"][0] = ptrs[i]
	}
}

func randomizeLevelNames(random *Random, rom []byte) {
	ndx := 0x21AC5
	for i := 0; i < len(TITLE_STRINGS); i++ {
		arr := TITLE_STRINGS[i]
		var str string
		if len(arr) > 1 {
			str = arr[random.NextInt(len(str))]
		} else {
			str = arr[0]
		}

		for j := 0; j < len(str); j++ {
			rom[ndx] = TEXT_MAPPING[uint16(str[j])]
		}

		rom[ndx-1] |= 0x80

	}
}

func charToTitleNum(chr uint8) byte {
	chars := map[uint8]byte{
		'@':  0x76, // clock
		'$':  0x2E, // coin
		'\'': 0x85,
		'"':  0x86,
		':':  0x78,
		' ':  0xFC}
	basechars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ.,*-!"
	for i := 0; i < len(basechars); i++ {
		chars[basechars[i]] = byte(i)
	}
	for k, _ := range chars {
		if k == chr {
			return chars[chr]
		}
	}
	return 0xFC
}

func writeToTitle(title string, color uint, rom []byte) {
	fmt.Println("writing ", title)
	title = strings.ToUpper(centerPad(title, 19))
	for i := 0; i < 19; i++ {
		var num = charToTitleNum(title[i])

		rom[0x2B6D7+i*2+0] = byte(num & 0xFF)
		rom[0x2B6D7+i*2+1] &= 0xE0
		rom[0x2B6D7+i*2+1] |= byte(byte((color << 2) | uint((num>>8)&0x3)))
		fmt.Printf("%c %d\n", title[i], num)
		printMD5(rom, "writeToTitle loop")
	}
}
