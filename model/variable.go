package model

import (
  "strconv"
)

type VariableType int

const (
	VariableReal VariableType = iota
	VariableInteger
	VariableBinary
	VariableArc
	VariableOther
)

func (t VariableType) String() string {
	switch t {
	case VariableReal:
		return "Real"
	case VariableInteger:
		return "Integer"
	case VariableBinary:
		return "Binary"
	case VariableArc:
		return "Arc"
	case VariableOther:
		return "Other"
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
