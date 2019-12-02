package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const input = "input.txt"

func main() {
	defer timeTrack(time.Now(), "main")

	file := loadFile(input)
	defer file.Close()

	ints := extractIntArr(file)
	var sum int
	c := make(chan int)
	for _, i := range ints {
		go func(i int) {
			c <- calculateFuel(i)
		}(i)
	}
	for range ints {
		sum += <-c
	}
	log.Printf("sum %d", sum)
}

func calculateFuel(i int) int {
	return i/3 - 2
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
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Panic("Error while converting to int")
		}
		ints = append(ints, i)
	}
	return ints
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
