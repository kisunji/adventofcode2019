package wirepanel

type Wire struct {
	id    int
	steps []Coordinates
}

func NewWire(id int, steps []Coordinates) *Wire {
	return &Wire{
		id:    id,
		steps: steps,
	}
}

func (w *Wire) AddStep(c Coordinates) {
	w.steps = append(w.steps, c)
}

func (w *Wire) GetLastCoords() Coordinates {
	return w.steps[len(w.steps)-1]
}

func (w *Wire) GetNumberOfSteps() int {
	return len(w.steps) - 1
}

func (w *Wire) GetId() int {
	return w.id
}
