package model

type ObjectiveSense int

const (
	ObjectiveMin ObjectiveSense = 0
	ObjectiveMax ObjectiveSense = 1
)

func (s ObjectiveSense) String() string {
	if s == ObjectiveMin {
		return "min"
	} else if s == ObjectiveMax{
		return "max"
	} else {
		return "unknown"
	}
}

type ObjectiveShape int

const (
	
)

func (s ObjectiveShape) String() string {
	return ""
}

type Objective struct {
	Name string
	Sense ObjectiveSense
	Shape ObjectiveShape
	Variables []Gradient
} 
