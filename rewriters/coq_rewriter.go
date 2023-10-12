package rewriters

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
	"github.com/tranduythanh/gocoqatoo/rewriters/rules"
)

type CoqRewriter struct {
	TextRewriter
}

func NewCoqRewriter(tr TextRewriter) *CoqRewriter {
	return &CoqRewriter{
		TextRewriter: tr,
	}
}

func (cr *CoqRewriter) GetTextVersion() string {
	textVersion := ""
	indentation := ""

	for i, p := range cr.InputsOutputs {
		input := p.Input
		output := p.Output
		var previousOutput *coq.Output
		assumptionsBeforeTactic := make(map[coq.Assumption]struct{})
		assumptionsAddedAfterTactic := make(map[coq.Assumption]struct{})
		for assumption := range output.Assumptions {
			assumptionsAddedAfterTactic[assumption] = struct{}{}
		}

		if i != 0 {
			previousOutput = cr.InputsOutputs[i-1].Output
			for assumption := range previousOutput.Assumptions {
				assumptionsBeforeTactic[assumption] = struct{}{}
			}
		}

		for k := range assumptionsBeforeTactic {
			delete(assumptionsAddedAfterTactic, k)
		}

		switch input.Typ {
		case coq.APPLY:
			textVersion += indentation
			textVersion += "(* "
			textVersion += rules.
				NewApply(cr.rewritingBundle).
				Apply(input, output, assumptionsBeforeTactic, assumptionsAddedAfterTactic, previousOutput)
			textVersion += " *) "
			textVersion += input.Value + "\n"

		case coq.ASSUMPTION:
			textVersion += indentation
			textVersion += "(* "
			// Assuming you have a Go function AssumptionApply() similar to the Java one
			textVersion += rules.
				NewAssumption(cr.rewritingBundle).
				Apply(nil, nil, nil, nil, nil)
			textVersion += " *) "
			textVersion += input.Value + "\n"

		case coq.BULLET:
			indentation = cr.UpdatedIndentationLevel(input)
			textVersion += indentation
			textVersion += input.Value
			textVersion += " (* "
			// Assuming a function cr.rewritingBundle[) to ]get string from bundle
			textVersion += fmt.Sprintf(cr.rewritingBundle["bullet"], "", output.Goal.Value)
			textVersion += " *)\n"
			indentation += "  "

		case coq.DESTRUCT:
			destructedObject := strings.TrimSpace(strings.Split(input.Value, " ")[1])
			textVersion += "(* "
			textVersion += fmt.Sprintf(cr.rewritingBundle["destruct"], destructedObject)
			textVersion += " *) "
			textVersion += input.Value + "\n"

		case coq.INTRO, coq.INTROS:
			textVersion += indentation
			textVersion += "(* "
			// Assuming a Go function IntrosApply() similar to the Java one
			textVersion += rules.
				NewIntros(cr.rewritingBundle).
				Apply(input, output, assumptionsBeforeTactic, assumptionsAddedAfterTactic, previousOutput)
			textVersion += " *) "
			textVersion += input.Value + "\n"

		case coq.INTUITION:
			textVersion += indentation
			textVersion += "(* "
			textVersion += fmt.Sprintf(cr.rewritingBundle["intuition"], previousOutput.Goal.Value)
			textVersion += " *)"
			textVersion += input.Value + "\n"

		case coq.INVERSION:
			textVersion += indentation
			textVersion += "(* "
			textVersion += rules.
				NewInversion(cr.rewritingBundle).
				Apply(input, output, assumptionsBeforeTactic, assumptionsAddedAfterTactic, nil)
			textVersion += " *) "
			textVersion += input.Value + "\n"

		case coq.LEMMA:
			textVersion += input.Value + "\n"

		case coq.OMEGA:
			textVersion += indentation
			textVersion += "(* "
			textVersion += cr.rewritingBundle["omega"]
			textVersion += " *)"
			textVersion += input.Value + "\n"

		case coq.PROOF:
			textVersion += input.Value + "\n"

		case coq.REFLEXIVITY:
			textVersion += indentation
			textVersion += "(* "
			textVersion += cr.rewritingBundle["reflexivity"]
			textVersion += " *)"
			textVersion += input.Value + "\n"

		case coq.SIMPL:
			textVersion += indentation
			textVersion += "(* "
			textVersion += fmt.Sprintf(cr.rewritingBundle["simpl"], previousOutput.Goal.Value, output.Goal.Value)
			textVersion += " *)"
			textVersion += input.Value + "\n"

		case coq.SPLIT:
			textVersion += indentation
			textVersion += input.Value + "\n"

		case coq.UNFOLD:
			textVersion += indentation
			unfoldedDefinition := strings.TrimSuffix(strings.Split(input.Value, " ")[1], ".")
			textVersion += "(* "
			textVersion += fmt.Sprintf(cr.rewritingBundle["unfold"], unfoldedDefinition, output.Goal.Value)
			textVersion += " *)"
			textVersion += input.Value + "\n"

		case coq.QED:
			textVersion += input.Value + "\n"

		default:
			textVersion += indentation
			textVersion += input.Value + "\n"
		}

	}

	textVersion = strings.ReplaceAll(textVersion, "<{[", "")
	textVersion = strings.ReplaceAll(textVersion, "}]>", "")
	return textVersion
}
