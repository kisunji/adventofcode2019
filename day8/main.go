package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	input = "input.txt"
	wide  = 25
	tall  = 6
)

func main() {
	file := loadFile(input)
	ints := extractIntArr(file)

	if len(ints)%(wide*tall) != 0 {
		log.Println("Invalid number of inputs")
	}

	layers := []Layer{}
	layer := NewLayer(wide, tall)

	for start, end := 0, wide; end <= len(ints); start, end = start+wide, end+wide {
		err := layer.Insert(ints[start:end])
		if err != nil {
			log.Panic("Insert failed!", err)
		}
		if layer.IsFull() {
			layers = append(layers, *layer)
			layer = NewLayer(wide, tall)
		}
	}

	leastZeroes := 999
	var layerLeastZeroes Layer

	for _, l := range layers {
		numZeroes := l.GetNumOfDigits(0)
		log.Printf("layer has %d zeroes", numZeroes)
		if numZeroes < leastZeroes {
			leastZeroes = numZeroes
			layerLeastZeroes = l
		}
	}

	num1Digits := layerLeastZeroes.GetNumOfDigits(1)
	num2Digits := layerLeastZeroes.GetNumOfDigits(2)
	log.Printf("result: %d", num1Digits*num2Digits)

	log.Printf("layers: %v", layers)

	// Prepare the final image by initializing to 2's
	finalImage := [tall][wide]int{}
	for i := 0; i < wide; i++ {
		for j := 0; j < tall; j++ {
			finalImage[j][i] = 2
		}
	}

	for i, m := range layers {
		log.Printf("Layer %d", i)
		for _, n := range m.data {
			log.Printf("%v\n", n)
		}
	}

	for i := 0; i < wide; i++ {
		for j := 0; j < tall; j++ {
			layerIndex := 0
			for finalImage[j][i] == 2 && layerIndex < len(layers) {
				finalImage[j][i] = layers[layerIndex].data[j][i]
				layerIndex++
			}
			log.Printf("%d,%d layer: %d", i, j, layerIndex-1)
		}
	}
	log.Printf("final layer")
	for _, m := range finalImage {
		log.Printf("%v\n", m)
	}
}

type Layer struct {
	width      int
	height     int
	data       [][]int
	currHeight int
}

var ErrInvalidInput = errors.New("Invalid input")
var ErrFullLayer = errors.New("Layer is already full")

func NewLayer(width, height int) *Layer {
	return &Layer{width: width, height: height}
}

func (l *Layer) IsFull() bool {
	return l.height == l.currHeight
}

func (l *Layer) Insert(input []int) error {
	if len(input) != l.width {
		return ErrInvalidInput
	}
	if l.IsFull() {
		return ErrFullLayer
	}
	l.data = append(l.data, input)
	l.currHeight++
	return nil
}

func (l *Layer) GetNumOfDigits(digit int) int {
	var count int
	for _, i := range l.data {
		for _, j := range i {
			if j == digit {
				count++
			}
		}
	}
	return count
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
		splitText := strings.Split(text, "")
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
