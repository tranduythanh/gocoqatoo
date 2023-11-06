package coq

import (
	"fmt"
	"strings"
)

type Assumption struct {
	Name string
	Typ  string
}

func (a *Assumption) ID() string {
	return fmt.Sprintf("%s:%s", a.Name, a.Typ)
}

func NewAssumption(value string) *Assumption {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return &Assumption{} // or handle error differently
	}
	return &Assumption{
		Name: strings.TrimSpace(parts[0]),
		Typ:  strings.TrimSpace(parts[1]),
	}
}

func (a *Assumption) TypeContainsSpaces() bool {
	return strings.Contains(a.Typ, " ")
}

func (a *Assumption) Equals(other *Assumption) bool {
	if other == nil {
		return false
	}
	return a.Name == other.Name && a.Typ == other.Typ
}

func (a *Assumption) HashCode() int {
	return hash(a.Typ) // assuming hash is a function that computes hash for string
}

// The Go language doesn't have a built-in hash function for strings.
// You might want to use a third-party package or a simple function like this:
func hash(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return h
}
