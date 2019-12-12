package main

import (
	"bufio"
	"github.com/kisunji/adventofcode2019/9/intcode"
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
	ints := extractInt64Arr(file)
	arr := append([]int64{}, ints...)

	computer := intcode.NewIntcodeComputer(arr)
	in := make(chan int64, 0)
	out := make(chan int64, 0)
	go computer.Compute(in, out)
	go func() { in <- 0 }()
	go func() { log.Println(<-out) }()
}


func loadFile(input string) *os.File {
	file, err := os.Open(input)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func extractInt64Arr(reader io.Reader) []int64 {
	var ints []int64
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		splitText := strings.Split(text, ",")
		for _, v := range splitText {
			i, err := strconv.ParseInt(v, 10, 64)
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
