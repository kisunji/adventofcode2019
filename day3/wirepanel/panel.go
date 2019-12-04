package wirepanel

import (
	"log"
	"math"
)

type lastVisited struct {
	wireId int
	steps  int
}

type Panel struct {
	grid           map[Coordinates]lastVisited
	closestByDist  int
	closestBySteps int
	central        Coordinates
}

func NewPanel() *Panel {
	return &Panel{
		grid:           map[Coordinates]lastVisited{},
		closestByDist:  math.MaxInt32,
		closestBySteps: math.MaxInt32,
		central:        Coordinates{},
	}
}

func (p *Panel) GetCentral() Coordinates {
	return p.central
}

func (p *Panel) GetClosestByDist() int {
	return p.closestByDist
}

func (p *Panel) GetClosestBySteps() int {
	return p.closestBySteps
}

func (p *Panel) RouteWireFull(w *Wire, insts []Instruction) {
	w.AddStep(p.GetCentral())
	for _, inst := range insts {
		p.RouteWireSingle(w, inst)
	}
}

func (p *Panel) RouteWireSingle(w *Wire, inst Instruction) {
	origin := w.GetLastCoords()
	log.Printf("wire[%d] from: %v instruction: %v", w.id, origin, inst)

	var current Coordinates
	for i := 1; i <= inst.distance; i++ {
		current = getCoordsDirectionOffset(origin, inst.direction, i)
		w.AddStep(current)
		p.CheckIntersection(w)
		p.SetLastVisited(w)
	}

	log.Printf("wire[%d] reached: %v", w.id, current)
}

func getCoordsDirectionOffset(origin Coordinates, direction string, offset int) Coordinates {
	x := origin.X
	y := origin.Y

	switch direction {
	case "U":
		return Coordinates{x, y + offset}
	case "R":
		return Coordinates{x + offset, y}
	case "D":
		return Coordinates{x, y - offset}
	case "L":
		return Coordinates{x - offset, y}
	default:
		log.Panicf("Invalid direction detected: %s", direction)
		return origin
	}
}

func (p *Panel) CheckIntersection(w *Wire) {
	coords := w.GetLastCoords()
	steps := w.GetNumberOfSteps()
	// If there is a non-zero lastVisited
	if lastVisited := p.grid[coords]; lastVisited.wireId != 0 && lastVisited.wireId != w.id {
		p.checkMinDistance(coords)
		p.checkMinSteps(coords, steps)
	}
}

func (p *Panel) checkMinDistance(c Coordinates) {
	newDistance := int(math.Abs(float64(c.X)) + math.Abs(float64(c.Y)))
	if newDistance < p.closestByDist {
		log.Printf("%v is new closestByDist! previous: %d, new: %d", c, p.closestByDist, newDistance)
		p.closestByDist = newDistance
	}
}

func (p *Panel) checkMinSteps(c Coordinates, steps int) {
	prevVisitor := p.grid[c]
	totalSteps := prevVisitor.steps + steps
	if totalSteps < p.closestBySteps {
		log.Printf("%v is new closestBySteps! previous: %d, new: %d", c, p.closestBySteps, totalSteps)
		p.closestBySteps = totalSteps
	}
}
func (p *Panel) SetLastVisited(w *Wire) {
	p.grid[w.GetLastCoords()] = lastVisited{
		wireId: w.GetId(),
		steps:  w.GetNumberOfSteps(),
	}
}
