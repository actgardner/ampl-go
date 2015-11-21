package model

import (
	"strconv"
)

type Gradient struct {
	Var Variable
	Coefficient float64
}

func (g Gradient) String() string {
	return strconv.FormatFloat(g.Coefficient, 'E', -1, 64) + " * " + g.Var.Name
}
