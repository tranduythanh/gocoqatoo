package rewriters

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
	"github.com/tranduythanh/gocoqatoo/rewriters/rules"
)

type TextRewriter struct {
	rewritingBundle         map[string]string // Assuming this is a map for now
	script                  string
	scriptWithUnfoldedAutos string
	InputsOutputs           []*coq.InputOutput
}

func NewTextRewriter(rewritingBundle map[string]string) *TextRewriter {
	return &TextRewriter{
		rewritingBundle: rewritingBundle,
	}
}

func (tr *TextRewriter) GenerateScriptWithUnfoldedAutos(inputsOutputs []*coq.InputOutput) string {
	scriptWithUnfoldedAutos := ""
	for _, p := range tr.InputsOutputs {
		i := p.Input
		o := p.Output
		if i.Typ == coq.AUTO {
			tacticsUsedByAuto := strings.Split(o.Value, "\n")
			for _, s := range tacticsUsedByAuto {
				if !strings.Contains(s, "(* info auto: *)") && !strings.Contains(s, "No more subgoals") {
					scriptWithUnfoldedAutos += strings.Replace(s, "simple ", "", -1) + "\n"
				}
			}
		} else {
			scriptWithUnfoldedAutos += i.Value + "\n"
		}
	}
	return scriptWithUnfoldedAutos
}

func (tr *TextRewriter) GetProofGlobalGoal() *coq.Goal {
	for _, p := range tr.InputsOutputs {
		i := p.Input
		o := p.Output
		if i.Typ == coq.LEMMA {
			return o.Goal
		}
	}
	return nil
}

func (tr *TextRewriter) UpdatedIndentationLevel(input *coq.Input) string {
	// Assumes that the input is of type BULLET
	if input.Typ == coq.BULLET {
		indentationLevel := len(strings.Split(input.Value, " ")[0]) // The bullet length determines the indentation level (e.g., - = 1, -- = 2, --- = 3)
		indentation := ""
		for i := 1; i <= indentationLevel; i++ {
			indentation += "  "
		}
		return indentation
	}
	return ""
}

func (tr *TextRewriter) GetTextVersion() string {
	textVersion := ""
	indentation := ""

	for i, p := range tr.InputsOutputs {
		input := p.Input
		output := p.Output
		var previousOutput *coq.Output
		assumptionsBeforeTactic := make(map[coq.Assumption]struct{})
		assumptionsAddedAfterTactic := make(map[coq.Assumption]struct{})

		for assumption := range output.Assumptions {
			assumptionsAddedAfterTactic[assumption] = struct{}{}
		}

		if i != 0 {
			previousOutput = tr.InputsOutputs[i-1].Output
			for assumption := range previousOutput.Assumptions {
				assumptionsBeforeTactic[assumption] = struct{}{}
			}
		}

		for assumption := range assumptionsBeforeTactic {
			delete(assumptionsAddedAfterTactic, assumption)
		}

		switch input.Typ {
		case coq.APPLY:
			textVersion += indentation
			textVersion += rules.
				NewApply(tr.rewritingBundle).
				Apply(input, output, assumptionsBeforeTactic, assumptionsAddedAfterTactic, previousOutput) + "\n"
		case coq.ASSUMPTION:
			textVersion += indentation
			textVersion += rules.
				NewAssumption(tr.rewritingBundle).
				Apply(nil, nil, nil, nil, nil) + "\n"
		case coq.BULLET:
			indentation = tr.UpdatedIndentationLevel(input)
			textVersion += indentation
			textVersion += fmt.Sprintf(tr.rewritingBundle["bullet"]+"\n", input.Value, output.Goal.Value)
			indentation += "  "
		case coq.DESTRUCT:
			textVersion += indentation
			textVersion += rules.
				NewDestruct(tr.rewritingBundle).
				Apply(input, nil, nil, nil, nil) + "\n"
		case coq.INTRO, coq.INTROS:
			textVersion += indentation
			textVersion += rules.
				NewIntros(tr.rewritingBundle).
				Apply(input, output, assumptionsBeforeTactic, assumptionsAddedAfterTactic, previousOutput) + "\n"
		case coq.LEMMA:
			textVersion += input.Value + "\n"
		case coq.OMEGA:
			textVersion += indentation
			textVersion += tr.rewritingBundle["omega"] + "\n"
		case coq.PROOF:
			textVersion += input.Value + "\n"
		case coq.REFLEXIVITY:
			textVersion += indentation
			textVersion += tr.rewritingBundle["reflexivity"] + "\n"
		case coq.SIMPL:
			textVersion += indentation
			textVersion += fmt.Sprintf(tr.rewritingBundle["simpl"]+"\n", previousOutput.Goal.Value, output.Goal.Value)
		case coq.UNFOLD:
			textVersion += indentation
			unfoldedDefinition := strings.Split(input.Value, " ")[1]
			unfoldedDefinition = strings.ReplaceAll(unfoldedDefinition, ".", "")
			textVersion += fmt.Sprintf(tr.rewritingBundle["unfold"]+"\n", unfoldedDefinition, output.Goal.Value)
		case coq.QED:
			textVersion += "Qed\n"
		}
	}

	return textVersion
}

func (tr *TextRewriter) Rewrite(proofScript string) {
	// formattedScript := tr.FormatScript(proofScript)
	// pp.Println(formattedScript)

	tr.ExtractInformation(proofScript)

	textVersion := tr.GetTextVersion()

	textVersion = strings.ReplaceAll(textVersion, "<[{", "")
	textVersion = strings.ReplaceAll(textVersion, "}]>", "")
	fmt.Println(textVersion)
}

func (tr *TextRewriter) FormatScript(proofScript string) string {
	formattedScript := proofScript

	// Step 1: Format proof so that there is one tactic/chain per line
	formattedScript = strings.ReplaceAll(formattedScript, ".", ".\n")

	// Step 2: Remove bullets
	lines := strings.Split(formattedScript, "\n")
	formattedScript = ""
	for _, line := range lines {
		l := strings.TrimSpace(line)
		for strings.HasPrefix(l, "-") || strings.HasPrefix(l, "*") || strings.HasPrefix(l, "+") || strings.HasPrefix(l, "{") || strings.HasPrefix(l, "}") {
			l = strings.TrimSpace(l[1:])
		}
		if l != "" {
			formattedScript += l + "\n"
		}
	}

	// Step 3: Execute formatted script to get the list of inputs/outputs
	// Assuming `Main.coqtop.execute` is a function returning a slice of InputOutput
	tr.InputsOutputs = coq.NewCoqtop(true).Execute(formattedScript)

	// Step 4: Build the proof tree
	// TODO: Clean this part
	var s []int
	bulletLevel := make(map[int]string)
	bulletsToAddAfter := make(map[int]string)
	bulletStr := ""
	for i, p := range tr.InputsOutputs {
		if p.Input.Value == "Qed." {
			break
		}
		if i == 0 {
			s = append(s, i)
		} else if i > 0 {
			previousNode := s[len(s)-1]
			s = s[:len(s)-1]

			previousPair := tr.InputsOutputs[i-1]
			numberOfSubgoalsBeforeTactic := previousPair.Output.GetNumberOfRemainingSubgoals()
			numberOfSubgoalsAfterTactic := p.Output.GetNumberOfRemainingSubgoals()

			addedSubgoals := numberOfSubgoalsAfterTactic - numberOfSubgoalsBeforeTactic
			for j := 0; j <= addedSubgoals; j++ {
				s = append(s, i)
			}
			if addedSubgoals > 0 {
				bulletStr += "-"
				bulletLevel[i] = bulletStr
			} else if addedSubgoals == 0 {
				if val, ok := bulletLevel[previousNode]; ok {
					bulletsToAddAfter[i] = val
				}
				s = append(s, i)
			} else if addedSubgoals < 0 {
				if val, ok := bulletLevel[previousNode]; ok {
					bulletsToAddAfter[i] = val
				}
				if len(s) > 0 {
					nextNodeId := s[len(s)-1]
					bulletStr = bulletLevel[nextNodeId]
				}
			}
		}
	}

	// Step 5: Insert bullets in inputsOutputs
	numberOfInputsInserted := 0
	for index, val := range bulletsToAddAfter {
		newInputOutput := &coq.InputOutput{
			Input:  &coq.Input{Value: val},
			Output: tr.InputsOutputs[index+numberOfInputsInserted-1].Output,
		}
		tr.InputsOutputs = append(tr.InputsOutputs[:index+numberOfInputsInserted], append([]*coq.InputOutput{newInputOutput}, tr.InputsOutputs[index+numberOfInputsInserted:]...)...)
		numberOfInputsInserted++
	}

	return formattedScript
}

func (tr *TextRewriter) ExtractInformation(proofScript string) {
	tr.script = proofScript
	if tr.InputsOutputs == nil {
		tr.InputsOutputs = coq.NewCoqtop(true).Execute(tr.script)
	}

	tr.scriptWithUnfoldedAutos = tr.GenerateScriptWithUnfoldedAutos(tr.InputsOutputs)
	if tr.scriptWithUnfoldedAutos != tr.script {
		coqtop := coq.NewCoqtop(true)
		tr.InputsOutputs = coqtop.Execute(tr.scriptWithUnfoldedAutos)
	}
}

func (r *TextRewriter) OutputProofTreeAsDot() {
	fmt.Println("---------------------------------------------")
	fmt.Println("|                Proof Tree                 |")
	fmt.Println("---------------------------------------------")
	fmt.Println("digraph {")

	stack := []int{}
	bulletLevel := make(map[int]string)
	bulletsToAddAfter := make(map[int]string)
	bulletStr := ""

	for i, p := range r.InputsOutputs {
		if p.Input.Value == "Qed." {
			break
		}
		if i == 0 {
			stack = append(stack, i)
		} else {
			previousNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // pop the last element

			previousPair := r.InputsOutputs[i-1]
			numberOfSubgoalsBeforeTactic := previousPair.Output.GetNumberOfRemainingSubgoals()
			numberOfSubgoalsAfterTactic := p.Output.GetNumberOfRemainingSubgoals()

			addedSubgoals := numberOfSubgoalsAfterTactic - numberOfSubgoalsBeforeTactic
			if addedSubgoals > 0 {
				for j := 0; j <= addedSubgoals; j++ {
					stack = append(stack, i)
				}
				fmt.Printf("\"%d. %s\" -> \"%d. %s\";\n", previousNode, r.InputsOutputs[previousNode].Input.Value, i, r.InputsOutputs[i].Input.Value)
				bulletStr += "-"
				bulletLevel[i] = bulletStr
			} else if addedSubgoals == 0 {
				if val, ok := bulletLevel[previousNode]; ok {
					bulletsToAddAfter[i] = val
				}
				fmt.Printf("\"%d. %s\" -> \"%d. %s\";\n", previousNode, r.InputsOutputs[previousNode].Input.Value, i, r.InputsOutputs[i].Input.Value)
				stack = append(stack, i)
			} else if addedSubgoals < 0 {
				if val, ok := bulletLevel[previousNode]; ok {
					bulletsToAddAfter[i] = val
				}
				fmt.Printf("\"%d. %s\" -> \"%d. %s\";\n", previousNode, r.InputsOutputs[previousNode].Input.Value, i, r.InputsOutputs[i].Input.Value)
				if len(stack) > 0 {
					nextNodeId := stack[len(stack)-1]
					bulletStr = bulletLevel[nextNodeId]
				}
			}
		}
	}
	fmt.Println("}")
}
