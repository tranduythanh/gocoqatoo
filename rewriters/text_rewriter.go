package rewriters

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
	"github.com/tranduythanh/gocoqatoo/rewriters/rules"
	"github.com/tranduythanh/gocoqatoo/stack"
	"github.com/yudai/pp"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type TextRewriter struct {
	rewritingBundle         map[string]string // Assuming this is a map for now
	script                  string
	scriptWithUnfoldedAutos string
	InputsOutputs           []*coq.InputOutput
	g                       *simple.DirectedGraph
}

func NewTextRewriter(rewritingBundle map[string]string) *TextRewriter {
	return &TextRewriter{
		rewritingBundle: rewritingBundle,
		g:               simple.NewDirectedGraph(),
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
	indentation := "\t-"

	for i, p := range tr.InputsOutputs {
		input := p.Input
		output := p.Output
		var previousOutput *coq.Output
		assumptionsBeforeTactic := make(map[string]*coq.Assumption)
		assumptionsAddedAfterTactic := make(map[string]*coq.Assumption)

		for aID, assumption := range output.Assumptions {
			assumptionsAddedAfterTactic[aID] = assumption
		}

		if i != 0 {
			previousOutput = tr.InputsOutputs[i-1].Output
			for aID, assumption := range previousOutput.Assumptions {
				assumptionsBeforeTactic[aID] = assumption
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
	tr.ExtractInformation(proofScript)

	textVersion := tr.GetTextVersion()

	textVersion = strings.ReplaceAll(textVersion, "<[{", "")
	textVersion = strings.ReplaceAll(textVersion, "}]>", "")
	fmt.Println("---------------------------------------------")
	fmt.Println("|               Text Version                |")
	fmt.Println("---------------------------------------------")
	fmt.Println(textVersion)
}

func (r *TextRewriter) OutputProofTreeAsDot() {
	fmt.Println("---------------------------------------------")
	fmt.Println("|                Proof Tree                 |")
	fmt.Println("---------------------------------------------")
	fmt.Println("digraph {")

	r.ParseGraph()
	txt := r.OutputGraphAsDot()
	fmt.Println(txt)
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

func (r *TextRewriter) OutputGraphAsDot() string {
	var dotBuilder strings.Builder
	dotBuilder.WriteString("digraph {\n")

	// Create a snapshot of nodes
	nodesSnapshot := make([]int64, 0)
	for nodes := r.g.Nodes(); nodes.Next(); {
		nodesSnapshot = append(nodesSnapshot, nodes.Node().ID())
	}

	pp.Println(nodesSnapshot)

	// Iterate over the snapshot of nodes
	for _, nodeID := range nodesSnapshot {
		curNode := r.InputsOutputs[nodeID]
		dotBuilder.WriteString(fmt.Sprintf("\t%d [label=\"%d. %s\"];\n", nodeID, nodeID, curNode.Input.Value))

		// Create a snapshot of edges for the current node
		edgesSnapshot := make([]graph.Edge, 0)
		for toNode := r.g.From(nodeID); toNode.Next(); {
			edgesSnapshot = append(edgesSnapshot, simple.Edge{
				F: simple.Node(nodeID),
				T: simple.Node(toNode.Node().ID()),
			})
		}

		// Iterate over the snapshot of edges
		for _, edge := range edgesSnapshot {
			targetID := edge.To().ID()
			dotBuilder.WriteString(fmt.Sprintf("\t%d -> %d;\n", nodeID, targetID))
		}
	}

	dotBuilder.WriteString("}\n")
	return dotBuilder.String()
}

func (r *TextRewriter) ParseGraph() {
	nodeMap := make(map[int]simple.Node)

	stack := stack.New[int]()

	for i, p := range r.InputsOutputs {
		if p.Input.Value == "Qed." {
			break
		}

		if i == 0 {
			stack.Push(i)
			continue
		}

		previousNode := stack.Pop()

		previousPair := r.InputsOutputs[i-1]
		numberOfSubgoalsBeforeTactic := previousPair.Output.GetNumberOfRemainingSubgoals()
		numberOfSubgoalsAfterTactic := p.Output.GetNumberOfRemainingSubgoals()

		addedSubgoals := numberOfSubgoalsAfterTactic - numberOfSubgoalsBeforeTactic

		if addedSubgoals > 0 {
			for j := 0; j <= addedSubgoals; j++ {
				stack.Push(i)
			}
			addEdge(r.g, nodeMap, previousNode, i)
			continue
		}

		if addedSubgoals == 0 {
			addEdge(r.g, nodeMap, previousNode, i)
			stack.Push(i)
			continue
		}

		if addedSubgoals < 0 {
			addEdge(r.g, nodeMap, previousNode, i)
			continue
		}
	}
}

func addEdge(g *simple.DirectedGraph, nodeMap map[int]simple.Node, src, dst int) {
	tryToAddNewNode(g, nodeMap, src, dst)
	g.SetEdge(g.NewEdge(nodeMap[src], nodeMap[dst]))
}

func tryToAddNewNode(g *simple.DirectedGraph, nodeMap map[int]simple.Node, nodes ...int) {
	for _, node := range nodes {
		if _, ok := nodeMap[node]; !ok {
			nodeMap[node] = simple.Node(node)
			g.AddNode(nodeMap[node])
		}
	}
}
