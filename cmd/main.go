package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tranduythanh/gocoqatoo/bundle"
	"github.com/tranduythanh/gocoqatoo/coq"
	"github.com/tranduythanh/gocoqatoo/rewriters"
)

var (
	inputFile  = flag.String("input", "", "File containing the Coq proof")
	language   = flag.String("language", "en", "Target language [en (default) | vi | fr]")
	mode       = flag.String("mode", "text", "Output mode [text (default) | coq | latex | dot]")
	debugMode  = flag.Bool("debug", false, "Display debugging information")
	locale     string
	bundleLang = map[string]string{}
)

func main() {
	flag.Parse()

	locale = *language

	switch locale {
	case "en":
		bundleLang = bundle.BundleEN
	case "fr":
		bundleLang = bundle.BundleFR
	case "vi":
		bundleLang = bundle.BundleVI
	default:
		fmt.Println("Unsupported language. Coqatoo currently supports: vi, en, fr.")
		os.Exit(0)
	}

	var debug = false

	if *debugMode {
		debug = true
	}

	coqtop := coq.NewCoqtop(debug)

	textRewriter := rewriters.NewTextRewriter(bundleLang)
	coqRewriter := rewriters.NewCoqRewriter(*textRewriter)
	latexRewriter := rewriters.NewLatexRewriter(*textRewriter)

	if *inputFile != "" {
		verifyFileExists(*inputFile)
		fileContents, _ := os.ReadFile(*inputFile)

		fmt.Println("---------------------------------------------")
		fmt.Println("|             Coq Version                   |")
		fmt.Println("---------------------------------------------")
		fmt.Println(string(fileContents))

		switch *mode {
		case "coq":
			fmt.Println("---------------------------------------------")
			fmt.Println("|                Coq Version                |")
			fmt.Println("---------------------------------------------")
			// Assuming you have a CoqRewriter function
			coqRewriter.Rewrite(string(fileContents))
		case "latex":
			fmt.Println("---------------------------------------------")
			fmt.Println("|              LaTeX Version                |")
			fmt.Println("---------------------------------------------")
			// Assuming you have a LatexRewriter function
			latexRewriter.Rewrite(string(fileContents))
		case "dot":
			// Assuming you have a TextRewriter function
			textRewriter.Rewrite(string(fileContents))
			// Assuming you have a outputProofTreeAsDot function
			textRewriter.OutputProofTreeAsDot()
		default:
			fmt.Println("---------------------------------------------")
			fmt.Println("|               Text Version                |")
			fmt.Println("---------------------------------------------")
			// Assuming you have a TextRewriter function
			textRewriter.Rewrite(string(fileContents))
		}

		coqtop.Stop()
	} else {
		flag.PrintDefaults()
	}
}

func verifyFileExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", filePath)
		os.Exit(1)
	}
}
