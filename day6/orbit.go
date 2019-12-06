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

	m := orbitMap{}

	var orbits int

	for _, s := range strs {
		inputs := strings.Split(s, ")")
		m[inputs[1]] = inputs[0]
	}

	for _, v := range m {
		orbits++
		centre, exists := m[v]
		for exists {
			orbits++
			centre, exists = m[centre]
		}
	}

	log.Printf("orbits :%d", orbits)

	youCentre, exists := m["YOU"]
	centre := youCentre
	youToCom := []string{}
	for exists {
		youToCom = append(youToCom, centre)
		centre, exists = m[centre]
	}
	sanCentre, exists := m["SAN"]
	centre = sanCentre
	sanToCom := []string{}
	for exists {
		sanToCom = append(sanToCom, centre)
		centre, exists = m[centre]
	}

	log.Printf("YOU->COM: %v", youToCom)
	log.Printf("SAN->COM: %v", sanToCom)

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

type orbitMap map[string]string

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
