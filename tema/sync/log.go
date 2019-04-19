package sync

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func (proc *Process) printlnOK(sync func(func()), text string) {
	proc.printf(sync, "%s\n", aurora.Green(text))
}

func (proc *Process) printlnSKIP(sync func(func()), text string) {
	proc.printf(sync, "%s\n", aurora.Yellow(text))
}

func (proc *Process) printlnERROR(sync func(func()), text string) {
	proc.printf(sync, "%s\n", aurora.Red(text))
}

func (proc *Process) printf(sync func(func()), template string, a ...interface{}) {
	proc.print(sync, fmt.Sprintf(template, a...))
}

func (proc *Process) println(sync func(func()), text string) {
	proc.print(sync, text+"\n")
}

func (proc *Process) print(sync func(func()), text string) {
	if !proc.quiet {
		if sync == nil {
			fmt.Print(text)
		} else {
			sync(func() { fmt.Print(text) })
		}
	}
}
