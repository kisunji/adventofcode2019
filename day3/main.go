package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"

	wp "github.com/kisunji/adventofcode2019/day3/wirepanel"
)

const (
	input = "input.txt"
	test1 = "test1.txt"
	test2 = "test2.txt"
)

func main() {
	defer timeTrack(time.Now(), "main")
	file := loadFile(input)
	//file := loadFile(test1)
	//file := loadFile(test2)
	defer file.Close()
	mtx := extractStringMatrix(file)

	p := wp.NewPanel()

	for i, v := range mtx {
		// Start Wire ID at 1 to avoid overlap with zero-value of lastVisited
		w := wp.NewWire(i+1)
		instructions := wp.ParseInstructions(v)
		p.RouteWireFull(w, instructions)
	}

	log.Printf("closestByDist: %d", p.GetClosestByDist())
	log.Printf("closestBySteps: %d", p.GetClosestBySteps())
}

func loadFile(input string) *os.File {
	file, err := os.Open(input)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func extractStringMatrix(reader io.Reader) [][]string {
	var m [][]string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var arr []string
		text := scanner.Text()
		splitText := strings.Split(text, ",")
		for _, v := range splitText {
			arr = append(arr, v)
		}
		m = append(m, arr)
	}
	return m
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
