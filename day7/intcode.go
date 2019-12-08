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

	phases := []int{0, 1, 2, 3, 4}
	permutations := [][]int{}
	for p := make([]int, len(phases)); p[0] < len(p); nextPerm(p) {
		permutations = append(permutations, getPerm(phases, p))
	}

	var highestSignal int

	for _, perm := range permutations {
		log.Printf("Permutation: %v", perm)

		inputChan := make(chan int)
		outputChan := make(chan int)

		inputA := []int{perm[0], 0}
		go func() {
			outputChan <- compute(ints, inputChan)
		}()
		for _, v := range inputA {
			inputChan <- v
		}
		outputA := <-outputChan

		inputB := []int{perm[1], outputA}
		go func() {
			outputChan <- compute(ints, inputChan)
		}()
		for _, v := range inputB {
			inputChan <- v
		}
		outputB := <-outputChan

		inputC := []int{perm[2], outputB}
		go func() {
			outputChan <- compute(ints, inputChan)
		}()
		for _, v := range inputC {
			inputChan <- v
		}
		outputC := <-outputChan

		inputD := []int{perm[3], outputC}
		go func() {
			outputChan <- compute(ints, inputChan)
		}()
		for _, v := range inputD {
			inputChan <- v
		}
		outputD := <-outputChan

		inputE := []int{perm[4], outputD}
		go func() {
			outputChan <- compute(ints, inputChan)
		}()
		for _, v := range inputE {
			inputChan <- v
		}
		outputE := <-outputChan

		log.Printf("output: %v", outputE)
		if outputE > highestSignal {
			highestSignal = outputE
		}
	}
	log.Printf("highestSignal: %d", highestSignal)
}

func stdin() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input: ")
	text, _ := reader.ReadString('\n')
	var err error
	strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		log.Panic(err)
	}
}

type instruction struct {
	opcode int
	modes  []int
}

func compute(instructions []int, c chan int) int {
	// Clone array
	arr := append([]int{}, instructions...)
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
	var output int
	for current.opcode != 99 {
		switch current.opcode {
		case 1:
			addOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 2:
			multOp(arr, pos, current.modes)
			pos += opParamMap[current.opcode] + 1
		case 3:
			inputOp(arr, pos, c)
			pos += opParamMap[current.opcode] + 1
		case 4:
			output = outputOp(arr, pos, current.modes)
			log.Printf("Output detected: %d", output)
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
	return output
}

func parseInstruction(i int, m map[int]int) *instruction {
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

func inputOp(arr []int, index int, c chan int) {
	log.Println("Waiting for input")
	input := <-c
	log.Printf("Got input: %d", input)
	target := arr[index+1]
	arr[target] = input
}

func outputOp(arr []int, index int, modes []int) int {
	target := arr[index+1]
	var output int
	if modes[0] == 0 {
		output = arr[target]
	} else {
		output = target
	}
	return output
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
