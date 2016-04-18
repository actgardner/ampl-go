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
	}
	return "Unknown"
}

type Variable struct {
	Name       string
	Type       VariableType
	lowerBound float64
	upperBound float64
	Index      int
}

/* Get the upper bound */
func (v Variable) UpperBound() float64 {
	return v.upperBound
}

/* Get the lower bound */
func (v Variable) LowerBound() float64 {
	return v.lowerBound
}

func (v Variable) String() string {
	str := "Name: " + v.Name
	str += " Type: " + v.Type.String()
	str += " Min: " + strconv.FormatFloat(v.LowerBound(), 'E', -1, 64)
	str += " Max: " + strconv.FormatFloat(v.UpperBound(), 'E', -1, 64)
	return str
}
