package wirepanel

import (
	"log"
	"strconv"
)

type Instruction struct {
	direction string
	distance  int
}

func ParseInstructions(s []string) []Instruction {
	instructions := []Instruction{}
	for _, instructionStr := range s {
		instruction := Instruction{
			direction: extractDirection(instructionStr),
			distance:  extractDistance(instructionStr),
		}
		instructions = append(instructions, instruction)
	}
	return instructions
}

func extractDirection(s string) string {
	return s[0:1]
}

func extractDistance(s string) int {
	distance, err := strconv.Atoi(s[1:])
	if err != nil {
		log.Panic(err)
	}
	return distance
}


