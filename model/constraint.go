package model

import (
	"strconv"
)

type ConstraintSense int

const (
	ConstraintGreaterThan ConstraintSense = iota
	ConstraintLessThan
	ConstraintEqualTo
	ConstraintRange
	ConstraintNonBinding
)

func (t ConstraintSense) String() string {
	switch t {
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
	Name      string
	Sense     ConstraintSense
	Shape     Shape
	Min       float64
	Max       float64
	Variables []Variable
	Index     int
	p         *Problem
}

/* Get the value of the constraint at a given point. Returns an error if AMPL is unable to compute the value */
func (c Constraint) Value(x []float64) (float64, error) {
	return c.p.conValue(c.Index, x)
}

/* Get the gradient of the constraint at a given point. Returns a nil slice and an error if AMPL is unable to compute the gradient */
func (c Constraint) Gradient(x []float64) ([]float64, error) {
	return c.p.conGrad(c.Index, x)
}

func (c Constraint) String() string {
	str := "Name: " + c.Name
	str += " Shape: " + c.Shape.String()
	str += " Sense: " + c.Sense.String()
	str += " ("
	switch c.Sense {
	case ConstraintGreaterThan:
		str += ">= " + strconv.FormatFloat(c.Min, 'E', -1, 64)
	case ConstraintLessThan:
		str += "<= " + strconv.FormatFloat(c.Max, 'E', -1, 64)
	case ConstraintEqualTo:
		str += "== " + strconv.FormatFloat(c.Min, 'E', -1, 64)
	case ConstraintRange:
		str += strconv.FormatFloat(c.Min, 'E', -1, 64) + "< <" + strconv.FormatFloat(c.Max, 'E', -1, 64)
	}
	str += ")"
	str += "\r\n"
	for _, v := range c.Variables {
		str += "\t" + v.String() + "\r\n"
	}
	return str
}

/* Returns a bool for whether or not the given value (as computed by `Value`) satisfies this constraint */
func (c Constraint) IsSatisfied(value float64) bool {
	switch c.Sense {
	case ConstraintGreaterThan:
		return value > c.Min-Featol
	case ConstraintLessThan:
		return value < c.Max+Featol
	case ConstraintEqualTo:
		return value > c.Min-Featol && value < c.Min+Featol
	case ConstraintRange:
		return value > c.Min-Featol && value < c.Max+Featol
	}
	return true
}
