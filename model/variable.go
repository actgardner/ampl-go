package model

import (
  "strconv"
)

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

func (v Variable) String() string {
	str := "Name: " + v.Name
	str += " Type: " + v.Type.String()
	str += " Min: " + strconv.FormatFloat(v.LowerBound, 'E', -1, 64)
	str += " Max: " + strconv.FormatFloat(v.UpperBound, 'E', -1, 64)
	return str
}
