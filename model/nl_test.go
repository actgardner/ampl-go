package model

import (
	"fmt"
	"testing"
)

const testModelFile = "diet1.nl"

func TestVariables(t *testing.T) {
	p := ProblemFromFile(testModelFile)
	vars := p.Variables()
	fmt.Printf("Vars: %v\n", vars)
}

func TestConstraints(t *testing.T) {
	p := ProblemFromFile(testModelFile)
	con := p.Constraints()
	fmt.Printf("Con: %v\n", con)
}

func TestObjectives(t *testing.T) {
	p := ProblemFromFile(testModelFile)
	obj := p.Objectives()
	fmt.Printf("Obj: %v\n", obj)
}
