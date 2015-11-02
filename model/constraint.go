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
}

