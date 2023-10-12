package rewriters

type Rewriter interface {
	Rewrite(proofScript string)
}
