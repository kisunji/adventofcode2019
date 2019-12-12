package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	input = "input.txt"
	test  = "test.txt"
	test2 = "test2.txt"
)

func main() {
	file := loadFile(input)
	//file := loadFile(test2)
	strs := extractStrArr(file)

	// Use map to track unidirectional adjacency
	// Allows us to trace path to root from ends of the tree
	// Key ORBITS value (e.g. A)B -> m[b] = a)
	type orbitMap map[string]string
	m := orbitMap{}

	for _, s := range strs {
		inputs := strings.Split(s, ")")
		m[inputs[1]] = inputs[0]
	}

	var orbits int

	// Range through all adjacencies, incrementing orbits until root is reached
	for _, v := range m {
		orbits++
		centre, exists := m[v]
		for exists {
			orbits++
			centre, exists = m[centre]
		}
	}

	log.Printf("orbits :%d", orbits)
	// Trace path from "YOU" to root
	youCentre, exists := m["YOU"]
	centre := youCentre
	youToCom := []string{}
	for exists {
		youToCom = append(youToCom, centre)
		centre, exists = m[centre]
	}
	// Trace path from "SAN" to root
	sanCentre, exists := m["SAN"]
	centre = sanCentre
	sanToCom := []string{}
	for exists {
		sanToCom = append(sanToCom, centre)
		centre, exists = m[centre]
	}

	log.Printf("YOU->COM: %v", youToCom)
	log.Printf("SAN->COM: %v", sanToCom)

	// Starting from root (COM), increment the two arrays
	// until a branching is detected. Count the leftover
	// orbits in each list
	var maxLength int
	if len(sanToCom) > len(youToCom) {
		maxLength = len(sanToCom)
	} else {
		maxLength = len(youToCom)
	}

	for i := 1; i < maxLength-1; i++ {
		a := len(sanToCom) - i
		b := len(youToCom) - i
		if sanToCom[a] != youToCom[b] {
			// By the time we detect a fork, we are already one node past it.
			// Add 2 to compensate on both branches
			log.Println(a + b + 2)
			return
		}
	}
}

func loadFile(input string) *os.File {
	file, err := os.Open(input)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func extractStrArr(reader io.Reader) []string {
	var strs []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		strs = append(strs, text)
	}
	return strs
}
