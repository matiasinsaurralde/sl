package main

import(
  "github.com/matiasinsaurralde/sl/parser"

  "fmt"
)

func main() {
  fmt.Println("Test Program\n")
  f, err := parser.ParseFile("ejemplos/holamundo.sl")

  if err != nil {
    panic(err)
  }

  fmt.Println("")
  fmt.Println("Ast.File:")

  fmt.Println(f)

  parser.Run( f )
}
