package ampl

/* A package to wrap the AMPL CLI and automate running commands */

import (
	"bufio"
	"fmt"
	"os/exec"
	"os"
	"strings"
	"time"
	"io"
	"io/ioutil"
)

const (
	// A file to read by default with the path to the AMPL executable 
	AMPLLocFile = "amplloc.txt"

	// Time to wait for an error from AMPL by default. Each command will take this long to return
	defaultTimeout = 500 * time.Millisecond
)

/* An instance of the AMPL binary */
type Runner struct {
	command *exec.Cmd
	stdout io.ReadCloser
	stdin io.WriteCloser
	output chan string
}

/* Look for a file called `amplloc.txt` and read the contents if it exists */
func GetAMPLLoc() (string, error) {
	amplLoc, err := ioutil.ReadFile(AMPLLocFile)
	// If the file is missing, return an empty string
	if os.IsNotExist(err) {
		return "", nil
	}
	return string(amplLoc), err
}

/* Present a prompt and try to get the AMPL location on stdin. Cache the location in `amplloc.txt` */
func PromptAMPLLoc() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("AMPL binary location: ")
	amplLoc, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	amplLoc = strings.TrimSpace(amplLoc)
	fmt.Printf("amplLoc: %q\n", amplLoc)
	ioutil.WriteFile(AMPLLocFile, []byte(amplLoc), os.ModePerm)
	return amplLoc, nil
}

/* Open a new runner using the AMPL binary at `amplLoc` */ 
func NewRunner(amplLoc string) (*Runner, error) {
	cmd := exec.Command(amplLoc)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	outputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(stdout)
		var output string
		for {
			l, err := reader.ReadString('\n')
			if err != nil {
				outputChan <- output
				break
			}
			// AMPL ends every error with an extra carriage return
			if l == "\r\r\n" {
				outputChan <- output
				output = ""	
			} else {
				output += string(l)
			}
		}
	}()
	return &Runner {
		command: cmd,
		stdout: stdout,
		stdin: stdin,
		output: outputChan,
	}, nil
}

/* Submit a command to the running AMPL binary, with the default timeout for a response. Returns an AMPLCommandError if the command fails, or nil if everything appears successful */
func (self *Runner) RunCommand (cmd string) (error) {
	return self.RunCommandTimeout(cmd, defaultTimeout)
}

/* Submit a command to the running AMPL binary. Waits for up to `timeout` for an error message from AMPL. Returns an AMPLCommandError
   if the command fails, or nil if there is no response. Note that if `timeout` is not long enough we may miss the error message. */
func (self *Runner) RunCommandTimeout (cmd string, timeout time.Duration) (error) {
	// Append a semicolon and newline to ensure the command is executed
	if cmd[len(cmd)-1] != ';' {
		cmd = cmd + ";"
	}
	_, err := self.stdin.Write([]byte(cmd+"\r\n"))

	if err != nil {
		return err
	}

	/* Unfortunately, the only feedback AMPL gives is when something goes wrong - listen for a configurable timeout after each command to see if it succeeded */	
	select {
	case errorString := <- self.output:
		return &AMPLCommandError {
			command: cmd,
			errorString: string(errorString),
		}
	case <- time.After(timeout):
		return nil
	}	
}

/* Returned when AMPL returns an error trying to run a command */
type AMPLCommandError struct {
  /* The command being run when the error happened */
  command string

  /* The text of the error */
  errorString string
}

func (self *AMPLCommandError) Error() string {
  return fmt.Sprintf("Error running %q - %v", self.command, self.errorString)
}
