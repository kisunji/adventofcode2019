package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

const input = "input.txt"

func main() {
	file := loadFile(input)
	asteroids := extractAsteroids(file)
	var optimal asteroid
	var mostLinesOfSight int
	for _, origin := range asteroids {
		set := map[float64]asteroid{}
		for _, other := range asteroids {
			if other == origin {
				continue
			}
			diffY := other.y - origin.y
			diffX := other.x - origin.x
			rad := math.Atan2(float64(diffY), float64(diffX))
			set[rad] = other
		}
		linesOfSight := len(set)

		log.Printf("lines of sight for %v: %v", origin, linesOfSight)
		if linesOfSight > mostLinesOfSight {
			optimal = origin
			mostLinesOfSight = linesOfSight
		}
	}
	log.Printf("optimal %v LOS: %v", optimal, mostLinesOfSight)

	// Part 2
	m := linkedAsteroidMap{}
	for _, other := range asteroids {
		if other == optimal {
			continue
		}
		diffY := other.y - optimal.y
		diffX := other.x - optimal.x
		rad := math.Atan2(float64(diffY), float64(diffX))
		m.add(rad, other)
	}

	var numVaporized int
	for len(m) > 0 {
		// First, extract the radians from the asteroid map and sort them
		rads := []float64{}
		for k, _ := range m {
			rads = append(rads, k)
		}
		sort.Float64s(rads)

		// Next, re-order the slice to start at -π/2 (vertical)
		var startIndex int
		for i, v := range rads {
			if v >= math.Atan2(-1, 0) {
				startIndex = i
				break
			}
		}
		modifiedRads := rads[startIndex:]
		modifiedRads = append(modifiedRads, rads[:startIndex]...)
		for _, key := range modifiedRads {
			numVaporized++
			log.Printf("Target %v: %v", numVaporized, m.popClosest(key, optimal))
		}
	}
}

// Note: y axis is flipped (goes downwards)
type asteroid struct {
	x int
	y int
}

type linkedAsteroidMap map[float64][]asteroid

func (l linkedAsteroidMap) add(rad float64, a asteroid) {
	if l[rad] == nil {
		l[rad] = []asteroid{a}
	} else {
		l[rad] = append(l[rad], a)
	}
}

func (l linkedAsteroidMap) popClosest(rad float64, source asteroid) asteroid {
	asteroids := l[rad]
	closest := source
	closestIndex := 0
	for i, a := range asteroids {
		if closest == source {
			closest = a
			closestIndex = i
		}
		if math.Abs(float64(a.x-source.x)) < math.Abs(float64(closest.x-source.x)) ||
			math.Abs(float64(a.y-source.y)) < math.Abs(float64(closest.y-source.y)) {
			closest = a
			closestIndex = i
		}
	}

	// Remove closest element from slice
	l[rad] = append(l[rad][:closestIndex], l[rad][closestIndex+1:]...) // If list is empty after removing, delete the key
	if len(l[rad]) == 0 {
		delete(l, rad)
	}

	return closest
}

func loadFile(input string) *os.File {
	file, err := os.Open(input)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func extractAsteroids(reader io.Reader) []asteroid {
	var result []asteroid
	scanner := bufio.NewScanner(reader)
	scanline := 0
	for scanner.Scan() {
		text := scanner.Text()
		splitText := strings.Split(text, "")
		for i, v := range splitText {
			if v == "#" {
				a := asteroid{i, scanline}
				result = append(result, a)
			}
		}
		scanline++
	}
	return result
}
