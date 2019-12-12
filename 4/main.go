package main

import (
	"log"
	"math"
	"time"
)

const (
	floor   = 367479
	ceiling = 893698
)

func main() {
	defer timeTrack(time.Now(), "main")
	p := &PasswordSolver{
		floor:   floor,
		ceiling: ceiling,
	}
	p.FindPasswords()
	log.Printf("Possible passwords: %v", p.GetPossiblePasswords())
	log.Printf("# of possible passwords: %d", p.GetNumPossiblePasswords())
	log.Printf("Possible passwords (v2): %v", p.GetPossiblePasswords2())
	log.Printf("# of possible passwords (v2): %d", p.GetNumPossiblePasswords2())
}

type PasswordSolver struct {
	floor              int
	ceiling            int
	possiblePasswords  []int
	possiblePasswords2 []int
}

func (p *PasswordSolver) FindPasswords() {
	for i := p.floor; i <= p.ceiling; i++ {
		i = p.scanInt(i)
	}
}
func (p *PasswordSolver) scanInt(i int) int {
	arr := intToArr(i)
	if len(arr) < 1 {
		log.Panic("Length of array is too small")
	}
	var hasAdjacentSame bool
	// Track number of occurrences for a given number
	// Multiple occurrences are guaranteed to be adjacent
	adjMap := map[int]int{}
	// Add first value
	adjMap[arr[0]]++
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] {
			hasAdjacentSame = true
		}
		// Increment current value by one
		adjMap[arr[i]]++
		if arr[i] < arr[i-1] {
			count := copyElToRight(i-1, arr)
			hasAdjacentSame = true
			// since we copied previous value {count} times,
			// we increment the previous value by {count}
			adjMap[arr[i-1]] += count
			// Exit loop because we know the following values were copied from previous
			break
		}
	}
	result := arrToInt(arr)
	if result > p.ceiling {
		return result
	}

	if hasAdjacentSame {
		p.AddPossiblePassword(result)
	}

	var hasUniquePair bool
	for _, v := range adjMap {
		if v == 2 {
			hasUniquePair = true
		}
	}

	if hasUniquePair {
		p.AddPossiblePassword2(result)
	}
	return result
}

func (p *PasswordSolver) AddPossiblePassword(i int) {
	p.possiblePasswords = append(p.possiblePasswords, i)
}

func (p *PasswordSolver) GetNumPossiblePasswords() int {
	return len(p.possiblePasswords)
}

func (p *PasswordSolver) GetPossiblePasswords() []int {
	return p.possiblePasswords
}

func (p *PasswordSolver) AddPossiblePassword2(i int) {
	p.possiblePasswords2 = append(p.possiblePasswords2, i)
}

func (p *PasswordSolver) GetNumPossiblePasswords2() int {
	return len(p.possiblePasswords2)
}

func (p *PasswordSolver) GetPossiblePasswords2() []int {
	return p.possiblePasswords2
}

// Returns number of overwritten elements
func copyElToRight(index int, arr []int) int {
	copyVal := arr[index]
	var counter int
	for i := index + 1; i < len(arr); i++ {
		arr[i] = copyVal
		counter++
	}
	return counter
}

func intToArr(i int) []int {
	if i < 0 {
		log.Panic("Cannot process negative integers")
	}
	arr := []int{}
	for i != 0 {
		arr = append([]int{i % 10}, arr...)
		i /= 10
	}
	return arr
}

func arrToInt(arr []int) int {
	var result int
	index := len(arr)
	for _, v := range arr {
		if v < 0 {
			log.Panic("Cannot process negative integers")
		}
		index--
		result += v * int(math.Pow10(index))
	}
	return result
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
