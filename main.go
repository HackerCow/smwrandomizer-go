package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	opt, err := parseFlags()
	check(err)

	dat, err := ioutil.ReadFile(opt.filename)
	check(err)

	err = checkMD5(dat)
	check(err)

	var seed int

	if opt.customSeed == "" {
		seed = int(rand.Int31())
	} else {
		seed64, err := strconv.ParseInt(opt.customSeed, 16, 64)
		seed = int(seed64)
		check(err)
	}

	fmt.Printf("Using seed %x\n", seed)
	fmt.Println("Randomizing...")
	randomizeRom(dat, seed, opt)
	fmt.Println("Done!")
}
