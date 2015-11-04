package model

type ConstraintShape int

const (
	NonLinearGeneralConstraint ConstraintShape = iota
	NonLinearNetworkConstraint
	LinearGeneralConstraint
	LinearNetworkConstraint
)

func (s ConstraintShape) String() string {
	switch (s) {
	case NonLinearGeneralConstraint:
		return "Non-Linear General"
	case NonLinearNetworkConstraint:
		return "Non-Linear Network"
	case LinearGeneralConstraint:
		return "Linear General"
	case LinearNetworkConstraint:
		return "Linear Network"
	}
	return "Unknown"
}

type ConstraintType int

const (
	ConstraintGreaterThan ConstraintType = iota
	ConstraintLessThan
	ConstraintEqualTo
	ConstraintRange
	ConstraintNonBinding
)

func (t ConstraintType) String() string {
	switch (t) {
	case ConstraintGreaterThan:
		return "Greater"
	case ConstraintLessThan:
		return "Less"
	case ConstraintEqualTo:
		return "Equals"
	case ConstraintRange:
		return "Range"
	case ConstraintNonBinding:
		return "Non-binding"
	}
	return "Unknown"
}

type Constraint struct {
 	Name string
	Shape ConstraintShape
	Type ConstraintType
	Min float64
	Max float64
	Variables []Gradient
	Index int
	p *Problem
}

/* Get the value of the constraint at a given point. Returns
   an error if AMPL is unable to compute the value */
func (c Constraint) Value(x []float64) (float64, error) {
	return c.p.conValue(c.Index, x)
}

/* Get the gradient of the constraint at a given point. Returns
   a nil slice and an error if AMPL is unable to compute the gradient */
func (c Constraint) Gradient(x []float64) ([]float64, error) {
	return c.p.conGrad(c.Index, x)
}

