package main

import (
	"fmt"
	"os"

	"github.com/matiasinsaurralde/sl/pkg/parser"
	"github.com/matiasinsaurralde/sl/pkg/runtime"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run cmd/sl/main.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	ast, err := parser.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	rt := runtime.NewRuntime()
	rt.RunFile(filename, ast)
}
