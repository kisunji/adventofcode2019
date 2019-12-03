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

const input = "input.txt"

func main() {
	defer timeTrack(time.Now(), "main")
	file := loadFile(input)
	mtx := extractStringMatrix(file)
	log.Println(mtx)
	p := &panel{grid: map[coords]visitedBy{}, closest: 0}
	for wire, instructions := range mtx {
		for _, i := range instructions {
			p.traverse(wire, i)
		}
	}
	log.Printf("%v", p.grid)
}

type panel struct {
	grid    map[coords]visitedBy
	closest int
}

type coords struct {
	x int
	y int
}

type visitedBy int

func (p *panel) traverse(wire int, input string) {
	direction := input[0:1]
	distance, err := strconv.Atoi(input[1:])
	if err != nil {
		log.Panic(err)
	}
	log.Printf("direction: %s, distance: %d", direction, distance)

	for i := 0; i <= distance; i++ {
		c := getCoordinates(direction, i)
		if v := p.grid[c]; int(v) != wire {
			log.Printf("intersection! wires %d and %d at %v", int(v), wire, c)
		}
		p.grid[c] = visitedBy(wire)
	}
}

func getCoordinates(direction string, i int) coords {
	switch direction {
	case "U":
		return coords{0, i}
	case "R":
		return coords{i, 0}
	case "D":
		return coords{0, -i}
	case "L":
		return coords{-i, 0}
	default:
		log.Panicf("Invalid direction detected: %s", direction)
		return coords{}
	}
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
