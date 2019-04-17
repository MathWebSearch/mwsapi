package sync

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func (proc *Process) printlnOK(text string) {
	fmt.Printf("%s\n", aurora.Green(text))
}

func (proc *Process) printlnERROR(text string) {
	fmt.Printf("%s\n", aurora.Red(text))
}

func (proc *Process) printf(template string, a ...interface{}) {
	proc.print(fmt.Sprintf(template, a...))
}

func (proc *Process) println(text string) {
	proc.print(text + "\n")
}

func (proc *Process) print(text string) {
	fmt.Print(text)
}
