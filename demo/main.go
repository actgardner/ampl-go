package main

import (
	runner "github.com/alanctgardner/ampl-go/runner"
	model "github.com/alanctgardner/ampl-go/model"
	"fmt"
)

/* Test commands to submit to the AMPL CLI - in this case
   use the diet model from the AMPL book and write an MPS file */
var amplCommands = []string {
	"model src/github.com/alanctgardner/ampl-go/demo/diet1.mod;",
	"data src/github.com/alanctgardner/ampl-go/demo/diet1.dat",
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
}
