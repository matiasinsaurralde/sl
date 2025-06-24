package test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/matiasinsaurralde/sl/pkg/parser"
	"github.com/matiasinsaurralde/sl/pkg/runtime"
)

// TestResult holds the result of running an SL program
type TestResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// RunSLProgram runs an SL program with custom input and captures output
func RunSLProgram(filename string, input string) (*TestResult, error) {
	// Open and parse the SL file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ast, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	// Create custom I/O streams
	stdin := strings.NewReader(input)
	var stdout, stderr bytes.Buffer

	// Create runtime with custom I/O
	rt := runtime.NewRuntimeWithIO(stdin, &stdout, &stderr)

	// Run the program
	rt.RunFile(filename, ast)

	return &TestResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0, // SL doesn't have exit codes yet
		Error:    nil,
	}, nil
}

// AssertOutput checks if the actual output matches expected output
func AssertOutput(t *testing.T, actual, expected string, testName string) {
	if actual != expected {
		t.Errorf("%s: unexpected output\nGot:  %q\nWant: %q", testName, actual, expected)
	}
}
