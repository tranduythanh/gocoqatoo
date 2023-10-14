package coq

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Coqtop struct {
	debugging bool
	process   *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	reader    *bufio.Reader
	writer    *bufio.Writer
}

func NewCoqtop(debugging bool) *Coqtop {
	coqtop := &Coqtop{
		debugging: debugging,
		process:   exec.Command("coqtop"),
	}

	var err error
	coqtop.stdin, err = coqtop.process.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return nil
	}

	coqtop.stdout, err = coqtop.process.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return nil
	}

	coqtop.stderr, err = coqtop.process.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return nil
	}

	if err := coqtop.process.Start(); err != nil {
		fmt.Println("Error starting process:", err)
		return nil
	}

	coqtop.reader = bufio.NewReader(coqtop.stdout)
	coqtop.writer = bufio.NewWriter(coqtop.stdin)

	coqtop.reader.ReadString('\n') // Ignore the first output of coqtop

	return coqtop
}

func (c *Coqtop) Execute(script string) []*InputOutput {
	scriptLines := strings.Split(script, "\n")
	var inputsOutputs []*InputOutput

	for _, input := range scriptLines {
		c.writer.WriteString(strings.ReplaceAll(input, "auto.", "info_auto.") + "\n")
		c.writer.Flush()

		output := ""

		timeoutTicker := time.NewTicker(time.Millisecond * 100)
		defer timeoutTicker.Stop()

	outerLoop:
		for {
			line, err := "", error(nil)
			readDone := make(chan bool)

			go func() {
				line, err = c.reader.ReadString('\n')
				readDone <- true
			}()

			select {
			case <-readDone:
				if err != nil && err != io.EOF {
					fmt.Println("Error reading:", err)
					break outerLoop
				}

				output += line
				if !strings.HasSuffix(line, "\n") {
					break outerLoop
				}

			case <-timeoutTicker.C:
				break outerLoop
			}
		}

		inputsOutputs = append(inputsOutputs, &InputOutput{
			Input:  NewInput(input),
			Output: NewOutput(output),
		})

		if c.debugging {
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			fmt.Println("Input:", strings.TrimSpace(input))
			fmt.Println("Output:\n", output)
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		}
	}

	return inputsOutputs
}

func (c *Coqtop) Stop() {
	if c.process.ProcessState == nil || !c.process.ProcessState.Exited() {
		c.process.Process.Kill()
	}
}
