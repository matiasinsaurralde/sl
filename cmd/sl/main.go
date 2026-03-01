// Command sl is the SL language interpreter.
// Usage: sl <file.sl> [args...]
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/matiasinsaurralde/sl/pkg/interpreter"
	"github.com/matiasinsaurralde/sl/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sl <file.sl> [args...]")
		os.Exit(1)
	}

	filename := os.Args[1]
	cmdArgs := os.Args[2:]

	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sl: cannot read %q: %v\n", filename, err)
		os.Exit(1)
	}

	prog, parseErrs := parser.Parse(string(src))
	if len(parseErrs) > 0 {
		for _, e := range parseErrs {
			fmt.Fprintf(os.Stderr, "sl: parse error: %v\n", e)
		}
		os.Exit(1)
	}

	interp := interpreter.New(cmdArgs)
	runErr := interp.Run(prog)
	if runErr != nil {
		var term *interpreter.TerminateError
		if errors.As(runErr, &term) {
			// terminar() with a message already printed it; exit 0
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "sl: %v\n", runErr)
		os.Exit(1)
	}
}
