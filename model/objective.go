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

/* Represents a single objective in the problem */
type Objective struct {
	Name string
	Sense ObjectiveSense
	Shape Shape
	Variables []Variable
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

func (o Objective) String() string {
	str := "Name: " + o.Name
	str += " Sense: " + o.Sense.String()
	str += " Shape: " + o.Shape.String()
	str += "\r\n"
	for _, v := range o.Variables {
		str += "\t" + v.String() + "\r\n"
	}
	return str
}
