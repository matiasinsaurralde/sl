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

	// x := "1 + 1 "

	// parser.ParseExpression(x)

	source := `
  var
  x : numerico
  b : numerico
  c = 1
  d : numerico
  e  = 2 + 2

  inicio
    a
    b
    c
  fin

  subrutina a(d)
  inicio
   hello_from_a
   bye_bye_from_a
  fin

  subrutina b(x,s)
  inicio
    hello_from_b
  fin

  subrutina complex(a,b) retorna numerico

  inicio
    hello_from_complex
  fin

  subrutina test_routine() retorna cadena

  inicio
    hello_from_test_routine
  fin
  `

	source = `
  var a = 1
  b = 2
  c = 1+1

  inicio
   a = 1
  fin
  `

	f, err := parser.Parse(source)

	fmt.Println("f =", f, " err = ", err)

	for i, node := range f.Nodes {
		// fmt.Println( "Node #", i, " = ", node)
		// fmt.Printf("%T\n", node)

		fmt.Print("Node #", i, " = ")
		switch v := node.(type) {
		case *Ast.GenericDeclaration:
			// fmt.Println("Node #", i, " = ", v, "(GenericDeclaration)" )
			fmt.Println(v, "(GenericDeclaration)")
			decl := node.(*Ast.GenericDeclaration)
			fmt.Println(decl, decl.Values)
		case *Ast.MainDeclaration:
			fmt.Println(v, "(MainDeclaration)")
		case *Ast.SubroutineDeclaration:
			fmt.Println(v, "(SubroutineDeclaration)")
		default:
			fmt.Println(v, "(other)")
		}
	}

	// expr2, err2 := goparser.ParseExpr("")

	// var node goast.ExprStmt
	// node = expr2

	/*
	  var b *goast.CallExpr
	  b = expr2.(*goast.CallExpr)

	  fmt.Println("1", b, *b)
	*/

	// fmt.Println("e",expr2, expr2.Pos(), expr2.End())

}
