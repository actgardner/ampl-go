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

/* Represents a single objective in the problem */
type Objective struct {
	Name string
	Sense ObjectiveSense
	Shape ObjectiveShape
	Variables []Gradient
	Index int
	p *Problem 
}

/* Compute the LHS of this objective at the given point */ 
func (o Objective) Value(x []float64) (float64, error) {
	return o.p.objValue(o.Index, x)	
}

/* Compute the gradient of this objective at the given point */
func (o Objective) Gradient(x []float64) ([]float64, error) {
	return o.p.objGrad(o.Index, x)
}
