package model

type VariableType int

const (
	VariableContinuousNonLinear VariableType = iota
	VariableIntegerNonLinear
	VariableLinearArc
	VariableOtherLinear
	VariableBinary
	VariableOtherInteger
)

func (t VariableType) String() string {
	switch t {
	case VariableContinuousNonLinear:
		return "Continuous Non-Linear"
	case VariableIntegerNonLinear:
		return "Integer Non-Linear"
	case VariableLinearArc:
		return "Linear Arc"
	case VariableOtherLinear:
		return "Other Linear"
	case VariableBinary:
		return "Binary"
	case VariableOtherInteger:
		return "Other Integer"
	}
	return "Unknown"
}

type Variable struct {
	Name string
	Type VariableType
	LowerBound float64
	UpperBound float64
	Index int
}
