package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	opt := options{"smw.sfc", "test_temp.sfc", "12345678", true, false, false,
		false, LEVEL_NAMES_MATCH_STAGE, false, BOWSER_DEFAULT, false,
		POWERUP_DEFAULT, false, false, false, false, false,
		false, false, false, false}

	defer os.Remove("test_temp.sfc")

	dat, err := ioutil.ReadFile(opt.filename)
	if err != nil {
		t.Error(err)
	}

	err = checkMD5(dat)
	if err != nil {
		t.Error(err)
	}

	randomizeRom(dat, 0x12345678, &opt)

	md5s := fmt.Sprintf("%x", md5.Sum(dat))

	if md5s != "1ce36ac1bab5322aaacdf25cb9f26d64" {
		t.Errorf("Expected 1ce36ac1bab5322aaacdf25cb9f26d64, got %s", md5s)
		t.Fail()
	}

}
