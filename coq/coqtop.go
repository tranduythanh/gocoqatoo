package coq

import (
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type Coqtop struct {
	debugging bool
}

func NewCoqtop(debugging bool) *Coqtop {
	coqtop := &Coqtop{
		debugging: debugging,
	}

	return coqtop
}

func (c *Coqtop) Execute(script string) []*InputOutput {
	output := c.execute(script)

	outputs := splitByTextLineNumbers(output)

	scriptLines := strings.Split(script, "\n")
	var inputsOutputs []*InputOutput

	i := 0
	preOutput := ""

	for _, input := range scriptLines {
		newInput := NewInput(input)
		newOutput := NewOutput("")
		switch input {
		case "Proof.", "Qed.":
			newOutput = NewOutput(preOutput)
		case "":
			continue
		default:
			newOutput = NewOutput(outputs[i])
			preOutput = outputs[i]
			i++
		}

		inputsOutputs = append(inputsOutputs, &InputOutput{
			Input:  newInput,
			Output: newOutput,
		})

	}

	return inputsOutputs
}

func (c *Coqtop) execute(script string) string {
	scriptLines := strings.Split(script, "\n")

	cmd := exec.Command("coqtop")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	output := ""

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// Start command first
	go func() {
		defer stdin.Close()
		for _, input := range scriptLines {
			io.WriteString(stdin, input+"\n")
		}
	}()

	// Read output in goroutine
	go func() {
		defer stdout.Close()
		stdOutBytes, _ := io.ReadAll(stdout)
		output = string(stdOutBytes)
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	return output
}

func splitByTextLineNumbers(text string) []string {
	text = strings.Trim(text, "\n\r\t ")

	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return []string{}
	}

	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\n\r\t ")
	}

	if strings.Contains(strings.ToUpper(lines[0]), "WELCOME TO COQ") {
		lines = lines[1:]
	}

	lineNumbers := getGoalLineNumbers(lines)

	var result []string
	for k, lineNumber := range lineNumbers {
		from := lineNumber
		to := len(lines)
		if k+1 < len(lineNumbers) {
			to = lineNumbers[k+1]
		}
		if lineNumber < len(lines) {
			txt := strings.Join(lines[from:to], "\n")
			result = append(result, txt)
		}
	}

	return result
}

func getGoalLineNumbers(lines []string) []int {
	var goalLines []int
	goalRegex := regexp.MustCompile(`^(\d+) goal(s)?$`)

	for i, line := range lines {
		if goalRegex.MatchString(line) ||
			strings.Contains(line, "No more goals.") {
			goalLines = append(goalLines, i)
		}
	}

	return goalLines

}
