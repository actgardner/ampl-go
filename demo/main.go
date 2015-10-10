package main

import (
	ampl "github.com/alanctgardner/ampl-go"
	"fmt"
)

/* Test commands to submit to the AMPL CLI - in this case
   use the diet model from the AMPL book and write an MPS file */
var amplCommands = []string {
	"model src/github.com/alanctgardner/ampl-go/demo/diet1.mod;",
	"data src/github.com/alanctgardner/ampl-go/demo/diet1.dat",
	"write mdiet1;",
}

func main() {
	/* Try to get the cached location of the AMPL binary from a file on disk */
	amplLoc, err := ampl.GetAMPLLoc()
	if err != nil {
		fmt.Printf("Error getting cached AMPL executable location: %v\n", err)
		return
	}

	/* If there was no cached location, prompt the user and cache what they enter - NB: if the user's entry is wrong, delete the cached file to be prompted again. */
	if amplLoc == "" {
		amplLoc, err = ampl.PromptAMPLLoc()
		if err != nil {
			fmt.Printf("Error getting AMPL executable location on stdin: %v\n", err)
			return
		}	
	}

	/* Run the AMPL binary */
	amplRunner, err := ampl.NewRunner(amplLoc)
	if err != nil {
		fmt.Printf("Error starting AMPL executable %q: %v\n", amplLoc, err)
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
}
