package intcode

import "log"

type instruction struct {
	opcode int
	modes  []int
}

type IntcodeComputer struct {
	instructions []int64
	paramMap     map[int]int
	position     int64
	relativeBase int64
}

func NewIntcodeComputer(instructions []int64) *IntcodeComputer {
	paramMap := map[int]int{
		1: 3,
		2: 3,
		3: 1,
		4: 1,
		5: 2,
		6: 2,
		7: 3,
		8: 3,
		9: 1,
	}
	return &IntcodeComputer{instructions, paramMap, 0, 0}
}

func (c *IntcodeComputer) Compute(in <-chan int64, out chan<- int64) {
	current := c.parseInstruction()
	for current.opcode != 99 {
		switch current.opcode {
		case 1:
			c.addOp(current.modes)
		case 2:
			c.multOp(current.modes)
		case 3:
			c.inputOp(current.modes, in)
		case 4:
			c.outputOp(current.modes, out)
		case 5:
			c.jumpIfTrueOp(current.modes)
		case 6:
			c.jumpIfFalseOp(current.modes)
		case 7:
			c.lessThanOp(current.modes)
		case 8:
			c.equalsOp(current.modes)
		case 9:
			c.relBaseOffsetOp(current.modes)
		default:
			log.Fatalln("Unsupported opcode!")
		}
		current = c.parseInstruction()
	}
}

func (c *IntcodeComputer) get(index int64) int64 {
	len64 := int64(len(c.instructions))
	if index >= len64 {
		// If index is out of range, extend instructions slice by the difference
		c.instructions = append(c.instructions, make([]int64, index-len64+1)...)
	}
	return c.instructions[index]
}

func (c *IntcodeComputer) set(index int64, value int64) {
	len64 := int64(len(c.instructions))
	if index >= len64 {
		// If index is out of range, extend instructions slice by the difference
		c.instructions = append(c.instructions, make([]int64, index-len64+1)...)
	}
	c.instructions[index] = value
}

func (c *IntcodeComputer) translateIndex(mode int, offset int) int64 {
	switch mode {
	case 0:
		return c.get(c.position + int64(offset))
	case 1:
		return c.position + int64(offset)
	case 2:
		relativeOffset := c.get(c.position + int64(offset))
		return relativeOffset + c.relativeBase
	default:
		log.Fatal("invalid index!")
		return 0
	}
}

func (c *IntcodeComputer) addOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	second := c.translateIndex(modes[1], 2)
	target := c.translateIndex(modes[2], 3)

	c.set(target, c.get(first)+c.get(second))

	c.position += 4
}

func (c *IntcodeComputer) multOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	second := c.translateIndex(modes[1], 2)
	target := c.translateIndex(modes[2], 3)

	c.set(target, c.get(first)*c.get(second))

	c.position += 4
}

func (c *IntcodeComputer) inputOp(modes []int, in <-chan int64) {
	log.Println(" Waiting for input")
	val := <-in
	log.Printf("Got input: %d", val)
	target := c.translateIndex(modes[0], 1)
	c.set(target, val)

	c.position += 2
}

func (c *IntcodeComputer) outputOp(modes []int, out chan<- int64) {
	target := c.translateIndex(modes[0], 1)
	out <- c.get(target)

	c.position += 2
}

func (c *IntcodeComputer) jumpIfTrueOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	target := c.translateIndex(modes[1], 2)
	if c.get(first) != 0 {
		c.position = c.get(target)
	} else {
		c.position += 3
	}
}

func (c *IntcodeComputer) jumpIfFalseOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	target := c.translateIndex(modes[1], 2)
	if c.get(first) == 0 {
		c.position = c.get(target)
	} else {
		c.position += 3
	}
}

func (c *IntcodeComputer) lessThanOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	second := c.translateIndex(modes[1], 2)
	target := c.translateIndex(modes[2], 3)
	if c.get(first) < c.get(second) {
		c.set(target, 1)
	} else {
		c.set(target, 0)
	}
	c.position += 4
}

func (c *IntcodeComputer) equalsOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	second := c.translateIndex(modes[1], 2)
	target := c.translateIndex(modes[2], 3)
	if c.get(first) == c.get(second) {
		c.set(target, 1)
	} else {
		c.set(target, 0)
	}
	c.position += 4
}

func (c *IntcodeComputer) relBaseOffsetOp(modes []int) {
	first := c.translateIndex(modes[0], 1)
	c.relativeBase += c.get(first)
	c.position += 2
}

func (c *IntcodeComputer) parseInstruction() *instruction {
	inst := c.instructions[c.position]
	opcode := int(inst) % 100
	modesRaw := intToArr(int(inst) / 100)
	modes := c.parseModes(modesRaw, c.paramMap[opcode])
	return &instruction{modes: modes, opcode: opcode}
}

func (c *IntcodeComputer) parseModes(arr []int, numParams int) []int {
	output := []int{}
	reverseInPlace(arr)
	output = append(output, arr...)
	for len(output) < numParams {
		output = append(output, 0)
	}
	return output
}

func reverseInPlace(a []int) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
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