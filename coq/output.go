package coq

import (
	"regexp"
	"strconv"
	"strings"
)

type Output struct {
	Value       string
	Assumptions map[string]*Assumption
	Goal        *Goal
}

func NewOutput(value string) *Output {
	value = strings.Trim(value, "\r\t\n ")
	return &Output{
		Value:       value,
		Assumptions: determineAssumptions(value),
		Goal:        determineGoal(value),
	}
}

func (o *Output) GetNumberOfRemainingSubgoals() int {
	if len(o.Value) == 0 {
		return 0
	}

	subgoalsStr := strings.Split(o.Value, " ")[0]

	subgoals, err := strconv.Atoi(subgoalsStr)
	if err != nil {
		subgoalsStr = strings.Split(o.Value, " ")[1]
		subgoals, err = strconv.Atoi(subgoalsStr)
		if err != nil {
			return 0
		}
	}

	return subgoals
}

func determineGoal(value string) *Goal {
	goals := parseGoals(value)
	if len(goals) > 0 {
		return NewGoal(strings.Join(goals, "\n"))
	}
	return NewGoal("")
}

func parseGoals(input string) []string {
	// Regular expression to match the goals
	// It looks for lines starting with "goal" followed by any text until the end of the line
	goalRegexp := regexp.MustCompile(`(?m)^goal\s\d\s*is:\s*(.*)$`)

	// Find all matches
	matches := goalRegexp.FindAllStringSubmatch(input, -1)

	// Initialize a slice to hold the goals
	var goals []string

	// Loop over the matches and extract the goals
	for _, match := range matches {
		if len(match) > 1 { // match[0] is the full match, match[1] is the first subgroup
			goals = append(goals, strings.TrimSpace(match[1]))
		}
	}

	// The first goal is not prefixed with "goal" in the input string, so we need to extract it separately
	// Split the input by lines
	lines := strings.Split(input, "\n")

	// Assuming the first goal is always on the line right after "============================"
	for i, line := range lines {
		if strings.Contains(line, "============================") && i+1 < len(lines) {
			firstGoal := strings.TrimSpace(lines[i+1])
			goals = append([]string{firstGoal}, goals...) // Prepend the first goal
			break
		}
	}

	return goals
}

func determineAssumptions(input string) map[string]*Assumption {
	// Initialize a map to hold the props and their types/definitions
	ret := map[string]*Assumption{}

	// Regular expression to match the props and their types/definitions
	// It looks for any pattern that starts with an identifier followed by ":" and the definition
	propRegexp := regexp.MustCompile(`(?m)^(\w+)\s*:\s*(.*)$`)

	// Find all matches
	matches := propRegexp.FindAllStringSubmatch(input, -1)

	// Loop over the matches and populate the map
	for _, match := range matches {
		if len(match) > 2 { // match[0] is the full match, match[1] is the name, match[2] is the typ
			name := strings.TrimSpace(match[1])
			typ := strings.TrimSpace(match[2])
			a := &Assumption{
				Name: name,
				Typ:  typ,
			}
			ret[a.ID()] = a
		}
	}

	return ret
}
