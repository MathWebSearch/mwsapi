package elasticsync

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// PrintlnOK prints text using "OK" semantics, optionally using syncronized output
func (proc *Process) PrintlnOK(sync func(func()), text string) {
	proc.Printf(sync, "%s\n", aurora.Green(text))
}

// PrintlnSKIP prints text using "SKIP" semantics, optionally using syncronized output
func (proc *Process) PrintlnSKIP(sync func(func()), text string) {
	proc.Printf(sync, "%s\n", aurora.Yellow(text))
}

// PrintlnERROR prints text using "ERROR" semantics, optionally using syncronized output
func (proc *Process) PrintlnERROR(sync func(func()), text string) {
	proc.Printf(sync, "%s\n", aurora.Red(text))
}

// Printf prints formatted text, optionally using syncronized output
func (proc *Process) Printf(sync func(func()), template string, a ...interface{}) {
	proc.Print(sync, fmt.Sprintf(template, a...))
}

// Println prints lined-delimited text, optionally using syncronized output
func (proc *Process) Println(sync func(func()), text string) {
	proc.Print(sync, text+"\n")
}

// Print prints text to stderr, optionally using syncronized output
func (proc *Process) Print(sync func(func()), text string) {
	if !proc.quiet {
		if sync == nil {
			os.Stderr.WriteString(text)
		} else {
			sync(func() { os.Stderr.WriteString(text) })
		}
	}
}
