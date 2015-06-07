package main

type room struct {
	out      *exit
	sublevel uint
	data     map[string][]byte
}
