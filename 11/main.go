package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kisunji/adventofcode2019/9/intcode"
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
	mode      int
}

func NewRobot() *Robot {
	origin := coords{0, 0}
	m := map[coords]int{}
	m[origin] = 1
	return &Robot{direction: UP, position: coords{0, 0}, path: m}
}

func (r *Robot) Run(arr []int64) {
	computer := intcode.NewIntcodeComputer(arr)
	inputChan := make(chan int64,1)
	outputChan := make(chan int64,1)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		computer.Compute(inputChan, outputChan)
		done <- true
		wg.Done()
	}()
	const (
		WRITE = 0
		MOVE  = 1
	)
	go func() {
		select {
		case <-done:
			return
		default:
			for {
				inputChan <- int64(r.read())
				x := <-outputChan
				log.Printf("Writing to pos(%v): %v", r.position, x)
				r.path[r.position] = int(x)
				r.mode = MOVE

				r.turn(int(<-outputChan))
				r.moveForward()
				r.mode = WRITE
			}
		}
		wg.Done()
	}()
	wg.Wait()
	log.Printf("Total tiles painted: %v", len(r.path))
}

func (r *Robot) read() int {
	r.mux.Lock()
	val := r.path[r.position]
	log.Printf("Read pos (%v):%v", r.position, r.path[r.position])
	r.mux.Unlock()
	return val
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
