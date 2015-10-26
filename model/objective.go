package model

type ObjectiveSense int

const (
	ObjectiveMin ObjectiveSense = 0
	ObjectiveMax ObjectiveSense = 1
)

type ObjectiveShape int

func (s ObjectiveSense) String() string {
	if s == ObjectiveMin {
		return "min"
	} else if s == ObjectiveMax{
		return "max"
	} else {
		return "unknown"
	}
}

type Objective struct {
	Name string
	Sense ObjectiveSense
	Shape ObjectiveShape
	Variables []Variable
} 
