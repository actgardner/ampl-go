package model

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const testModelFile = "diet1.nl"

/* Get a list of variables in the diet problem */
func TestDietVariables(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	vars := p.Variables()
	assert.Equal(len(vars), 9, "Number of variables")
	assert.Equal(vars[0], Variable{"_svar[1]", VariableInteger, 0, 11, 0}, "variable 1")
	assert.Equal(vars[1], Variable{"_svar[2]", VariableInteger, 0, 10, 1}, "variable 2")
	assert.Equal(vars[2], Variable{"_svar[3]", VariableInteger, 0, 8, 2}, "variable 3")
	assert.Equal(vars[3], Variable{"_svar[4]", VariableInteger, 0, 9, 3}, "variable 4")
	assert.Equal(vars[4], Variable{"_svar[5]", VariableInteger, 0, 8, 4}, "variable 5")
	assert.Equal(vars[5], Variable{"_svar[6]", VariableInteger, 0, 14, 5}, "variable 6")
	assert.Equal(vars[6], Variable{"_svar[7]", VariableInteger, 0, 13, 6}, "variable 7")
	assert.Equal(vars[7], Variable{"_svar[8]", VariableInteger, 0, 31, 7}, "variable 8")
	assert.Equal(vars[8], Variable{"_svar[9]", VariableInteger, 0, 18, 8}, "variable 9")
}

/* Get the list of constraints in the diet problem */
func TestDietConstraints(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	cons := p.Constraints()
	vars := p.Variables()
	assert.Equal(len(cons), 7, "number of constraints")
	assert.Equal(cons[0], Constraint{"_scon[1]", ConstraintGreaterThan, Linear, 2000, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 0, p}, "Contraint 1")

	assert.Equal(cons[1], Constraint{"_scon[2]", ConstraintRange, Linear, 350, 375, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 1, p}, "Contraint 2")

	assert.Equal(cons[2], Constraint{"_scon[3]", ConstraintGreaterThan, Linear, 55, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 2, p}, "Contraint 3")

	assert.Equal(cons[3], Constraint{"_scon[4]", ConstraintGreaterThan, Linear, 100, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[6], vars[7], vars[8]}, 3, p}, "Contraint 4")

	assert.Equal(cons[4], Constraint{"_scon[5]", ConstraintGreaterThan, Linear, 100, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[4], vars[5], vars[7], vars[8]}, 4, p}, "Contraint 5")

	assert.Equal(cons[5], Constraint{"_scon[6]", ConstraintGreaterThan, Linear, 100, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[6], vars[7], vars[8]}, 5, p}, "Contraint 6")

	assert.Equal(cons[6], Constraint{"_scon[7]", ConstraintGreaterThan, Linear, 100, math.Inf(1), []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[8]}, 6, p}, "Contraint 7")
}

/* Get the list of objectives in the diet problem */
func TestDietObjectives(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	objs := p.Objectives()
	vars := p.Variables()
	assert.Equal(len(objs), 8, "number of objectives")
	assert.Equal(objs[0], Objective{"_sobj[1]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 0, p})
	assert.Equal(objs[1], Objective{"_sobj[2]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 1, p})
	assert.Equal(objs[2], Objective{"_sobj[3]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 2, p})
	assert.Equal(objs[3], Objective{"_sobj[4]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[7], vars[8]}, 3, p})
	assert.Equal(objs[4], Objective{"_sobj[5]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[6], vars[7], vars[8]}, 4, p})
	assert.Equal(objs[5], Objective{"_sobj[6]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[4], vars[5], vars[7], vars[8]}, 5, p})
	assert.Equal(objs[6], Objective{"_sobj[7]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[6], vars[7], vars[8]}, 6, p})
	assert.Equal(objs[7], Objective{"_sobj[8]", ObjectiveMin, Linear, []Variable{vars[0], vars[1], vars[2], vars[3], vars[4], vars[5], vars[6], vars[8]}, 7, p})
}

/* Get the value of an objective in the diet problem */
func TestDietObjVal(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	objs := p.Objectives()
	val, err := objs[0].Value([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(val, float64(2.19), "Objective value")
}

/* Get the gradient of an objective in the diet problem */
func TestDietObjGradient(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	objs := p.Objectives()
	grad, err := objs[0].Gradient([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(grad, []float64{1.84, 2.19, 1.84, 1.44, 2.29, 0.77, 1.29, 0.6, 0.72}, "Objective gradient")
}

/* Get the value of a constraint in the diet problem */
func TestDietConVal(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	cons := p.Constraints()
	val, err := cons[0].Value([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(val, float64(370), "Constraint value")
}

/* Get the gradient of a constraint in the diet problem */
func TestDietConGradient(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	cons := p.Constraints()
	grad, err := cons[0].Gradient([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(grad, []float64{510, 370, 500, 370, 400, 220, 345, 110, 80}, "Constraint gradient")
}

/* Get all of the constraint values for the diet problem */
func TestDietConstraintValues(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	conVals, err := p.ConstraintValues([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(conVals, []float64{370, 35, 24, 15, 10, 20, 20}, "Constraint values")
}

/* Get the Jacobian of the constraint values for the diet problem */
func TestDietConstraintJacobian(t *testing.T) {
	assert := assert.New(t)
	p := ProblemFromFile(testModelFile)
	conVals, err := p.ConstraintJacobian([]float64{0, 1, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(err, "No error")
	assert.Equal(conVals, []float64{510, 34, 28, 15, 6, 30, 20, 370, 35, 24, 15, 10, 20, 20, 500, 42, 25, 6, 2, 25, 20, 370, 38, 14, 2, 15, 10, 400, 42, 31, 8, 15, 15, 8, 220, 26, 3, 15, 2, 345, 27, 15, 4, 20, 15, 110, 12, 9, 10, 4, 30, 80, 20, 1, 2, 120, 2, 2, 0, 0, 0, 0, 0}, "Constraint values")
}
