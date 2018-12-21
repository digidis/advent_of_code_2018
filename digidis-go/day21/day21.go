package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var instructions = map[string]int{
	"addr": 0, "addi": 1, "mulr": 2, "muli": 3,
	"banr": 4, "bani": 5, "borr": 6, "bori": 7,
	"setr": 8, "seti": 9, "gtir": 10, "gtri": 11,
	"gtrr": 12, "eqir": 13, "eqri": 14, "eqrr": 15,
}

var (
	values          = make(map[int]int)
	first, firstOps int
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Done in %v\n", time.Since(start))
	}()

	var (
		code [][]int
		ipr  int
	)

	data, _ := ioutil.ReadFile("test.txt")
	rows := strings.Split(string(data), "\n")

	for _, r := range rows {
		c := strings.Fields(r)
		if c[0] == "#ip" {
			ipr = atoi(c[1])
			continue
		}
		code = append(code, []int{
			instructions[c[0]], atoi(c[1]), atoi(c[2]), atoi(c[3]),
		})
	}

	run(code, ipr, 0)

	fmt.Printf("\n%d different values got it to halt\n", len(values))
	fmt.Printf("First one: %v\n", first)
	m := 0
	ops := 0
	for k, v := range values {
		if v > ops {
			m = k
			ops = v
		}
	}
	fmt.Printf("Most operations: %v (%d operations)\n", m, ops)
}

func run(code [][]int, ipr int, r0 int) (int, int) {
	reg := []int{r0, 0, 0, 0, 0, 0, 0}
	ip := 0
	ops := 0
	l := len(code)

	for {
		reg[ipr] = ip
		c := code[ip]
		switch c[0] {
		case 0:
			reg[c[3]] = reg[c[1]] + reg[c[2]]
		case 1:
			reg[c[3]] = reg[c[1]] + c[2]
		case 2:
			reg[c[3]] = reg[c[1]] * reg[c[2]]
		case 3:
			reg[c[3]] = reg[c[1]] * c[2]
		case 4:
			reg[c[3]] = reg[c[1]] & reg[c[2]]
		case 5:
			reg[c[3]] = reg[c[1]] & c[2]
		case 6:
			reg[c[3]] = reg[c[1]] | reg[c[2]]
		case 7: // bori
			reg[c[3]] = reg[c[1]] | c[2]
		case 8:
			reg[c[3]] = reg[c[1]]
		case 9: // seti
			reg[c[3]] = c[1]
		case 10:
			if c[1] > reg[c[2]] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		case 11:
			if reg[c[1]] > c[2] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		case 12:
			if reg[c[1]] > reg[c[2]] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		case 13:
			if c[1] == reg[c[2]] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		case 14:
			if reg[c[1]] == c[2] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		case 15:
			//------
			if first == 0 {
				first = reg[3]
				firstOps = ops
			}
			if values[reg[3]] != 0 {
				break
			}
			values[reg[3]] = ops
			fmt.Printf(".")
			//------
			if reg[c[1]] == reg[c[2]] {
				reg[c[3]] = 1
			} else {
				reg[c[3]] = 0
			}
		}

		ip = reg[ipr] + 1
		if ip >= l {
			break
		}
		ops++
	}
	return reg[0], ops
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
