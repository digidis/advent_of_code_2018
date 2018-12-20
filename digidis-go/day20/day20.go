package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"time"
)

var (
	doors = make(map[string]int)
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Done in %v\n", time.Since(start))
	}()

	data, _ := ioutil.ReadFile("input.txt")
	parse(data[1:len(data)-1], 0, 0, 0)
	m := 0
	var mk []byte
	more1k := 0
	for k, v := range doors {
		if v >= 1000 {
			more1k++
		}
		if v > m {
			m = v
			mk = []byte(k)
		}
	}
	fmt.Printf("Part1 : Room %d,%d is %v doors away\n", binary.LittleEndian.Uint16(mk[0:2]), binary.LittleEndian.Uint16(mk[2:]), m)
	fmt.Printf("Part2 : %d rooms\n", more1k)
}

func parse(r []byte, x, y int16, steps int) int {
	var (
		ox, oy, osteps = x, y, steps
		key            = make([]byte, 4)
	)
	for i := 0; i < len(r); i++ {
		switch r[i] {
		case 'N':
			y -= 1
			steps++
		case 'S':
			y += 1
			steps++
		case 'E':
			x += 1
			steps++
		case 'W':
			x -= 1
			steps++
		case '(':
			i += parse(r[i+1:], x, y, steps)
			i++
		case ')':
			return i
		case '|':
			x = ox
			y = oy
			steps = osteps
		}
		key[0] = byte(x & 0xff)
		key[1] = byte(x >> 8)
		key[2] = byte(y & 0xff)
		key[3] = byte(y >> 8)
		if v, ok := doors[string(key)]; !ok || v > steps {
			doors[string(key)] = steps
		}
	}
	return 0
}
