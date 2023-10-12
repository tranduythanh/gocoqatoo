package bundle

var BundleFR = map[string]string{
	"apply":          "Selon l'hypothèse <[{%s}]>, nous savons que <[{%s}]> est vrai lorsque <[{%s}]> sont vrais.",
	"assumption":     "Vrai, car c'est une de nos hypothèses.",
	"bullet":         "%s Cas <[{%s}]>:",
	"destruct":       "Considérons les différents cas possible de <[{%s}]>.",
	"intros.given":   "Soit <[{%s}]>.",
	"intros.suppose": "Supposons <[{%s}]>.",
	"intros.goal":    "Montrons <[{%s}]>.",
	"intuition":      "Intuitivement, on sait que <[{%s}]> est vrai.",
	"inversion":      "Par inversion sur <[{%s}]>, on obtient que %s est également vrai.",
	"omega":          "En utilisant l'algorithme Omega de William Pugh's, on peut conclure que c'est vrai.",
	"reflexivity":    "True.",
	"simpl":          "En simplifiant <[{%s}]>, on obtient que c'est équivalent à <[{%s}]>.",
	"unfold":         "En remplaçant <[{%s}]> par sa définition, on obtient <[{%s}]>.",
}
