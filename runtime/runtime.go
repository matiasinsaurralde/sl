package SL

import (
	"github.com/matiasinsaurralde/sl/ast"
	"github.com/matiasinsaurralde/sl/token"

	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Runtime struct{}

func NewRuntime(filename string) (*Runtime, error) {
	var err error
	runtime := Runtime{}
	_, err = ioutil.ReadFile(filename)
	return &runtime, err
}

func (runtime *Runtime) Start() {
	log.Println("Starting")
}

func Call(call *Ast.CallExpression) {
	fmt.Println("Call:", call.Function, "(function)\n")

	var node Ast.Node
	node = call.Args[0]
	// var l Ast.BasicLiteral
	l := *node.(*Ast.BasicLiteral)
	v := l.Value
	tok := token.Lookup(call.Function)

	switch tok {
	case token.PRINT:

		v = strings.Replace(v, "\"", "", -1)
		fmt.Println(v)
	}
}

func Run(f *Ast.File) {

	fmt.Println("\nRunning...\n")

	// Declarations:

	for _, d := range f.Scope.Declarations {
		var declaration *Ast.GenericDeclaration
		declaration = d.(*Ast.GenericDeclaration)

		fmt.Println("Declaring:", declaration)
	}

	fmt.Println("\nMain...\n")

	// for _, n := range f.Scope.Nodes {
	var Main *Ast.RoutineLiteral
	Main = f.Scope.Nodes[0].(*Ast.RoutineLiteral)

	// fmt.Println(*Main.Body)
	bs := *Main.Body

	vars := make(map[string]interface{}, 0)
	fmt.Println("vars", vars)

	for _, stPointer := range bs.List {
		var CallEx *Ast.CallExpression
		CallEx = stPointer.(*Ast.CallExpression)

		fmt.Println("Evaluate:", CallEx, "\n")

		Call(CallEx)
	}
	// fmt.Println( "Node", n )
	// }
}
