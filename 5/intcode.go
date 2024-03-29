package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	input = "input.txt"
)

func main() {
	defer timeTrack(time.Now(), "main")
	file := loadFile(input)
	defer file.Close()
	ints := extractIntArr(file)
	_ = computeInput(ints)
}

type instruction struct {
	opcode int
	modes  []int
}

func computeInput(input []int) []int {
	// Clone array
	arr := append([]int{}, input...)
	pos := 0
	opParamMap := map[int]int{
		1: 3,
		2: 3,
		3: 1,
		4: 1,
		5: 2,
		6: 2,
		7: 3,
		8: 3,
	}
	current := parseInstruction(arr[pos], opParamMap)
	for current.opcode != 99 {
		switch current.opcode {
		case 1:
			addOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 2:
			multOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 3:
			inputOp(arr, pos)
			pos += opParamMap[current.opcode] + 1
		case 4:
			outputOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 5:
			pos = jumpIfTrueOp(arr, pos, current.modes)
		case 6:
			pos = jumpIfFalseOp(arr, pos, current.modes)
		case 7:
			lessThanOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 8:
			equalsOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		default:
			log.Fatalln("Unsupported opcode!")
		}
		current = parseInstruction(arr[pos], opParamMap)
	}
	return arr
}

func parseInstruction(i int, m map[int]int) *instruction {
	log.Printf("instruction :%d", i)
	opcode := i % 100
	modesRaw := intToArr(i / 100)
	modes := parseModes(modesRaw, m[opcode])
	return &instruction{modes: modes, opcode: opcode}
}

func parseModes(arr []int, numParams int) []int {
	output := []int{}
	reverseInPlace(arr)
	output = append(output, arr...)
	for len(output) < numParams {
		output = append(output, 0)
	}
	return output
}

func addOp(arr []int, index int, modes []int) {
	var first int
	var second int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		second = arr[index+2]
	} else {
		second = index + 2
	}
	target := arr[index+3]
	arr[target] = arr[first] + arr[second]
}

func multOp(arr []int, index int, modes []int) {
	var first int
	var second int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		second = arr[index+2]
	} else {
		second = index + 2
	}
	target := arr[index+3]
	arr[target] = arr[first] * arr[second]
}

func inputOp(arr []int, index int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input: ")
	text, _ := reader.ReadString('\n')
	target := arr[index+1]
	var err error
	arr[target], err = strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		log.Panic(err)
	}
}

func outputOp(arr []int, index int, modes []int) {
	target := arr[index+1]
	var output int
	if modes[0] == 0 {
		output = arr[target]
	} else {
		output = target
	}
	fmt.Printf("Output: %d\n", output)
}

func jumpIfTrueOp(arr []int, index int, modes []int) int {
	var first int
	var target int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		target = arr[index+2]
	} else {
		target = index + 2
	}
	if arr[first] != 0 {
		return arr[target]
	}

	return index + 3
}

func jumpIfFalseOp(arr []int, index int, modes []int) int {
	var first int
	var target int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		target = arr[index+2]
	} else {
		target = index + 2
	}
	if arr[first] == 0 {
		return arr[target]
	}

	return index + 3
}

func lessThanOp(arr []int, index int, modes []int) {
	var first int
	var second int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		second = arr[index+2]
	} else {
		second = index + 2
	}
	target := arr[index+3]
	if arr[first] < arr[second] {
		arr[target] = 1
	} else {
		arr[target] = 0
	}
}

func equalsOp(arr []int, index int, modes []int) {
	var first int
	var second int
	if modes[0] == 0 {
		first = arr[index+1]
	} else {
		first = index + 1
	}
	if modes[1] == 0 {
		second = arr[index+2]
	} else {
		second = index + 2
	}
	target := arr[index+3]
	if arr[first] == arr[second] {
		arr[target] = 1
	} else {
		arr[target] = 0
	}
}

func reverseInPlace(a []int) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func intToArr(i int) []int {
	if i < 0 {
		log.Panic("Cannot process negative integers")
	}
	arr := []int{}
	for i != 0 {
		arr = append([]int{i % 10}, arr...)
		i /= 10
	}
	return arr
}

func loadFile(input string) *os.File {
	file, err := os.Open(input)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func extractIntArr(reader io.Reader) []int {
	var ints []int
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		splitText := strings.Split(text, ",")
		for _, v := range splitText {
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Panic(err)
			}
			ints = append(ints, i)
		}
	}
	return ints
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
