package main

import (
	"fmt"
	"math/rand"
	"time"
	runner "github.com/alanctgardner/ampl-go/runner"
	model "github.com/alanctgardner/ampl-go/model"
)

/* Test commands to submit to the AMPL CLI - in this case
   use the diet model from the AMPL book and write an MPS file */
var amplCommands = []string {
	"model comprehensive;",
	"option auxfiles rc;",
	"write gdiet1;",
}

func main() {
	/* Try to get the cached location of the AMPL binary from a file on disk */
	amplLoc, err := runner.GetAMPLLoc()
	if err != nil {
		fmt.Printf("Error getting cached AMPL executable location: %v\n", err)
		return
	}

	/* If there was no cached location, prompt the user and cache what they enter */
	if amplLoc == "" {
		amplLoc, err = runner.PromptAMPLLoc()
		if err != nil {
			fmt.Printf("Error getting AMPL executable location on stdin: %v\n", err)
			return
		}	
	}

	/* Run the AMPL binary */
	amplRunner, err := runner.NewRunner(amplLoc)
	if err != nil {
		fmt.Printf("Error starting AMPL executable %q: %v\n", amplLoc, err)
		runner.ClearAMPLLoc()
		return
	}

	/* Run each of the commands in order. If there's an error, print it out and quit. Continuing after an error will not work. */ 
	for _, cmd := range amplCommands {
		fmt.Printf("Running command: %q\n", cmd)
		err := amplRunner.RunCommand(cmd)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	/* Stop the AMPL process */
	amplRunner.Stop()

	/* Load the .nl file */
	p := model.ProblemFromFile("diet1.nl")

	/* Print all the variables in the problem */
	fmt.Printf("\nVariables\n---\n")
	for _, v := range p.Variables() {
		fmt.Printf("Variable %v - %s\n", v.Index, v)
	}

	/* Print all the constraints in the problem */
	fmt.Printf("\nConstraints\n---\n")
	for _, c := range p.Constraints() {
		fmt.Printf("Constraint %v - %s\n", c.Index, c)
	}

	/* Print all the objectives in the problem */
	fmt.Printf("\nObjectives\n---\n")
	for _, o := range p.Objectives() {
		fmt.Printf("Objective %v - %s\n", o.Index, o)
	}

	/* Get random values within the bounds of every variable */
	fmt.Printf("\nRandom variable values:\n---\n")
	vars := make([]float64, len(p.Variables())) 
	rand.Seed(time.Now().Unix())
	for i, v := range p.Variables() {
		vars[i] = (rand.Float64() * (v.UpperBound - v.LowerBound)) + v.LowerBound	
		fmt.Printf("\t%v = %v\n", v.Name, vars[i])
	}

	/* Evaluate every constraint at the randomly chosen point */
	fmt.Printf("\nConstraint values:\n---\n")
	for _, c := range p.Constraints() {
		conVal, err := c.Value(vars)
		if err != nil {
			fmt.Printf("Error evaluating constraint %v: %v\n", c.Name, err)
			return
		}
		fmt.Printf("\t%v = %v (satisfied: %v)\n", c.Name, conVal, c.IsSatisfied(conVal))		
	}

	/* Evaluate every objective at the random point */
	fmt.Printf("\nObjective values:\n---\n")
	for _, o := range p.Objectives() {
		objVal, err := o.Value(vars)
		if err != nil {
			fmt.Printf("Error evaluating objective %v: %v\n", o.Name, err)
			return
		}
		fmt.Printf("\t%v = %v\n", o.Name, objVal)		
	}

	/* Compute the gradient of every constraint at the randomly chosen point */
	fmt.Printf("\nConstraint gradients:\n---\n")
	for _, c := range p.Constraints() {
		conGrad, err := c.Gradient(vars)
		if err != nil {
			fmt.Printf("Error evaluating constraint gradient %v: %v\n", c.Name, err)
			return
		}
		fmt.Printf("\t%v = %v\n", c.Name, conGrad)	
	}

	/* Compute the gradient of every objective at the random point */
	fmt.Printf("\nObjective gradients:\n---\n")
	for _, o := range p.Objectives() {
		objGrad, err := o.Gradient(vars)
		if err != nil {
			fmt.Printf("Error evaluating objective gradient %v: %v\n", o.Name, err)
			return
		}
		fmt.Printf("\t%v = %v\n", o.Name, objGrad)
	}
}
