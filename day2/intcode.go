package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	input  = "input.txt"
	target = 19690720
)

func main() {
	defer timeTrack(time.Now(), "main")
	file := loadFile(input)
	defer file.Close()
	ints := extractIntArr(file)
	//ints[1] = 12
	//ints[2] = 2
	//result := computeInput(ints)
	//log.Println(result)

	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			ints[1] = i
			ints[2] = j
			result := computeInput(ints)
			if result[0] == target {
				log.Printf("noun: %d, verb %d", i, j)
				log.Printf("answer: %d", 100*i+j)
				return
			}
		}
	}
	log.Printf("%d not found", target)
}

func computeInput(arr []int) []int {
	clone := append([]int{}, arr...)
	index := 0
	opcode := clone[index]
	for opcode != 99 {
		switch opcode {
		case 1:
			addOp(clone, index)
		case 2:
			multOp(clone, index)
		default:
			log.Fatalln("Unsupported opcode!")
		}
		index += 4
		opcode = clone[index]
	}
	return clone
}

func addOp(arr []int, index int) {
	i := arr[index+1]
	j := arr[index+2]
	target := arr[index+3]
	arr[target] = arr[i] + arr[j]
}

func multOp(arr []int, opIndex int) {
	i := arr[opIndex+1]
	j := arr[opIndex+2]
	target := arr[opIndex+3]
	arr[target] = arr[i] * arr[j]
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
