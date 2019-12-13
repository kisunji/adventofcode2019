package main

import (
	"bufio"
	"github.com/kisunji/adventofcode2019/9/intcode"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	input = "input.txt"
)

const (
	UP    = "up"
	DOWN  = "down"
	LEFT  = "left"
	RIGHT = "right"
)

func main() {
	defer timeTrack(time.Now(), "main")
	file := loadFile(input)
	defer file.Close()
	ints := extractInt64Arr(file)
	arr := append([]int64{}, ints...)
	r := NewRobot()
	r.Run(arr)
}

type Robot struct {
	direction string
	position  coords
	path      map[coords]int
	mux       sync.Mutex
}

func NewRobot() *Robot {
	origin := coords{0, 0}
	m := map[coords]int{}
	m[origin] = 1
	return &Robot{direction: UP, position: coords{0, 0}, path: m}
}

func (r *Robot) Run(arr []int64) {
	computer := intcode.NewIntcodeComputer(arr)
	inputChan := make(chan int64, 1)
	outputChan := make(chan int64, 1)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		computer.Compute(inputChan, outputChan)
		done <- true
		wg.Done()
	}()
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
			default:
				// Following instructions should be blocking to maintain order
				// Robot sends read input to inputChan
				inputChan <- int64(r.read())
				// First val from outputChan is the writeVal
				r.write(int(<-outputChan))
				// Second val from outputChan is the command to turn
				r.turn(int(<-outputChan))
				r.moveForward()
			}
		}
	}()
	wg.Wait()
	log.Printf("Total tiles painted: %v", len(r.path))   // Iterate through coords in the map, finding the boundaries
	var minX, maxX, minY, maxY int
	for key, _ := range r.path {
		if key.x <= minX {
			minX = key.x
		} else if key.x >= maxX {
			maxX = key.x
		}
		if key.y <= minY {
			minY = key.y
		} else if key.y >= maxY {
			maxY = key.y
		}
	}
	offsetX := 0 - minX
	offsetY := 0 - minY
	width := maxX - maxY + 1
	height := maxY - minY + 1
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}
	for k, v := range r.path {
		log.Println(k)
		grid[k.y+offsetY][k.x+offsetX] = v
	}
	// Image is flipped vertically
	reverse(grid)
	for _, v := range grid {
		log.Println(v)
	}
}

func (r *Robot) read() int {
	val := r.path[r.position]
	log.Printf("Read pos (%v):%v", r.position, r.path[r.position])
	return val
}

func (r *Robot) write(writeVal int) {
	log.Printf("Writing to pos(%v): %v", r.position, writeVal)
	r.path[r.position] = writeVal
}

func (r *Robot) turn(command int) {
	log.Printf("turning from %v", r.direction)
	log.Printf("command: %d", command)
	switch r.direction {
	case UP:
		if command == 0 {
			r.direction = LEFT
		} else {
			r.direction = RIGHT
		}
	case RIGHT:
		if command == 0 {
			r.direction = UP
		} else {
			r.direction = DOWN
		}
	case DOWN:
		if command == 0 {
			r.direction = RIGHT
		} else {
			r.direction = LEFT
		}
	case LEFT:
		if command == 0 {
			r.direction = DOWN
		} else {
			r.direction = UP
		}
	}
	log.Printf("turned to %v", r.direction)
}

func (r *Robot) moveForward() {
	log.Printf("moving from %v", r.position)
	switch r.direction {
	case UP:
		r.position.y++
	case RIGHT:
		r.position.x++
	case DOWN:
		r.position.y--
	case LEFT:
		r.position.x--
	}
	log.Printf("moved to %v", r.position)
}

type coords struct {
	x int
	y int
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

func reverse(a [][]int) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}