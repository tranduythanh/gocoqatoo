package coq

import "strings"

type Input struct {
	Value string
	Typ   InputType
}

func NewInput(value string) *Input {
	trimmedValue := value // TODO: Remove existing comments from input
	typ := determineType(trimmedValue)
	return &Input{Value: trimmedValue, Typ: typ}
}

func determineType(value string) InputType {
	switch {
	case strings.HasPrefix(value, "Abort"):
		return ABORT
	case strings.HasPrefix(value, "apply"):
		return APPLY
	case strings.HasPrefix(value, "assumption"):
		return ASSUMPTION
	case strings.HasPrefix(value, "auto"):
		return AUTO
	case strings.HasPrefix(value, "-"), strings.HasPrefix(value, "+"), strings.HasPrefix(value, "*"):
		return BULLET // TODO: Add support for { } bullets
	case strings.HasPrefix(value, "destruct"):
		return DESTRUCT
	case strings.HasPrefix(value, "intros"):
		return INTROS
	case strings.HasPrefix(value, "intro"):
		return INTRO
	case strings.HasPrefix(value, "intuition"):
		return INTUITION
	case strings.HasPrefix(value, "inversion"):
		return INVERSION
	case strings.HasPrefix(value, "Lemma"):
		return LEMMA
	case strings.HasPrefix(value, "Proof"):
		return PROOF
	case strings.HasPrefix(value, "reflexivity"):
		return REFLEXIVITY
	case strings.HasPrefix(value, "simpl"):
		return SIMPL
	case strings.HasPrefix(value, "split"):
		return SPLIT
	case strings.HasPrefix(value, "unfold"):
		return UNFOLD
	case strings.HasPrefix(value, "Qed"):
		return QED
	default:
		return UNKNOWN
	}
}
