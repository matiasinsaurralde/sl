package main

import(
  "github.com/matiasinsaurralde/sl/parser"

  "fmt"
)

func main() {
  fmt.Println("test_ast")
  f, err := parser.ParseFile("ejemplos/holamundo.sl")

  if err != nil {
    panic(err)
  }

  fmt.Println(f)
}
