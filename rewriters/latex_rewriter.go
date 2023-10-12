package rewriters

import "strings"

type LatexRewriter struct {
	TextRewriter
}

func NewLatexRewriter(tr TextRewriter) *LatexRewriter {
	return &LatexRewriter{
		TextRewriter: tr,
	}
}

func (lr *LatexRewriter) GetTextVersion() string {
	latexVersion := lr.TextRewriter.GetTextVersion()
	latexVersion = strings.ReplaceAll(latexVersion, "<[{", "$")
	latexVersion = strings.ReplaceAll(latexVersion, "}]>", "$")
	latexVersion = strings.ReplaceAll(latexVersion, "Lemma", "\\begin{lemma}")
	latexVersion = strings.ReplaceAll(latexVersion, "Proof.", "\\end{lemma}\\begin{proof}")
	latexVersion = strings.ReplaceAll(latexVersion, "Qed", "\\end{proof}")
	latexVersion = strings.ReplaceAll(latexVersion, "and", "~and~")
	latexVersion = strings.ReplaceAll(latexVersion, "/\\", "\\land")
	latexVersion = strings.ReplaceAll(latexVersion, "<->", "\\Leftrightarrow")
	latexVersion = strings.ReplaceAll(latexVersion, "->", "\\Rightarrow")
	latexVersion = strings.ReplaceAll(latexVersion, "<-", "\\Leftarrow")
	latexVersion = strings.ReplaceAll(latexVersion, "forall", "\\forall")
	latexVersion = strings.ReplaceAll(latexVersion, "\n", "\\\\\n")
	latexVersion = strings.ReplaceAll(latexVersion, "  ", "\\hspace{5mm}")
	latexVersion = strings.ReplaceAll(latexVersion, "\\hspace{5mm}\\\\", "\\\\")
	latexVersion = strings.ReplaceAll(latexVersion, "\\hspace{5mm}\\\\", "\\\\")

	return latexVersion
}
