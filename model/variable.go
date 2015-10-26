package model

type VariableType int

const (
	VariableNonLinear VariableType = iota
	VariableNetworkLinear
	VariableOtherLinear
	VariableLinearBinary
	VariableInteger
)

func (t VariableType) String() string {
	switch t {
	case VariableNonLinear:
		return "Non-Linear"
	case VariableNetworkLinear:
		return "Network Linear"
	case VariableOtherLinear:
		return "Other Linear"
	case VariableLinearBinary:
		return "Binary"
	case VariableInteger:
		return "Integer"
	}
	return "Unknown"
}

type Variable struct {
	Name string
	Type VariableType
	LowerBound float64
	UpperBound float64
}
