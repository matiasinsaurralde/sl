package main

import (
	"github.com/matiasinsaurralde/sl/ast"
	"github.com/matiasinsaurralde/sl/parser"

	// goparser "go/parser"
	// goast "go/ast"

	"fmt"
)

func main() {
	fmt.Println("Test Program\n")

	source := `1+abc`
	e := parser.Eval(source, nil)

	fmt.Println(e)

	var be *Ast.BinaryExpression
	be = e.(*Ast.BinaryExpression)

	// var x, y *Ast.BinaryExpression
	// x = be.X.(*Ast.BinaryExpression)
	// y = be.Y.(*Ast.BinaryExpression)

	fmt.Println("BE\nX = ", be.X, "Y = ", be.Y)
	fmt.Println("OPER = ", be.Operator)

}
