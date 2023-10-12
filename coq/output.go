package coq

import (
	"strconv"
	"strings"
)

type Output struct {
	Value       string
	Assumptions map[Assumption]struct{}
	Goal        *Goal
}

func NewOutput(value string) *Output {
	return &Output{
		Value:       value,
		Assumptions: determineAssumptions(value),
		Goal:        determineGoal(value),
	}
}

func (o *Output) GetNumberOfRemainingSubgoals() int {
	subgoalsStr := strings.Split(o.Value, " ")[0]
	subgoals, err := strconv.Atoi(subgoalsStr)
	if err != nil {
		return 0
	}
	return subgoals
}

func determineGoal(value string) *Goal {
	t := strings.Split(value, "============================\n ")
	if len(t) > 1 {
		goalStr := strings.Split(t[1], "\n")[0]
		return NewGoal(strings.TrimSpace(goalStr))
	}
	return NewGoal("")
}

func determineAssumptions(value string) map[Assumption]struct{} {
	assumptions := make(map[Assumption]struct{})

	t := splitBySeparator(value)
	if len(t) <= 1 {
		return assumptions
	}

	assumpStrs := extractAssumptions(t[0])
	for _, s := range assumpStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		assumptionParts := splitAssumption(s)
		if len(assumptionParts) != 2 {
			continue
		}

		addAssumptions(assumptionParts[0], assumptionParts[1], assumptions)
	}

	return assumptions
}

func splitBySeparator(value string) []string {
	return strings.Split(value, "============================\n ")
}

func extractAssumptions(value string) []string {
	t := strings.Split(value, "\n  \n ")
	if len(t) > 1 {
		return strings.Split(t[1], "\n ")
	}
	return []string{}
}

func splitAssumption(s string) []string {
	return strings.Split(s, " : ")
}

func addAssumptions(namesStr, typeStr string, assumptions map[Assumption]struct{}) {
	names := strings.Split(namesStr, ", ")
	for _, name := range names {
		assumptions[*NewAssumption(name + " : " + typeStr)] = struct{}{}
	}
}
