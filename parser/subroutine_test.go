package parser_test

import(
  "github.com/matiasinsaurralde/sl/parser"
  "github.com/matiasinsaurralde/sl/ast"

  "testing"
)

func testSubroutineDeclaration( f *Ast.File, t *testing.T, subroutineName string) {

  var subroutineDeclaration *Ast.SubroutineDeclaration
  declaration := *(&f.Nodes[0])

  if len(f.Nodes) == 0 {
    t.Error("Ast.File Nodes is empty")
  }

  if len(f.Nodes) > 1 {
    t.Error("Ast.File Nodes is larger than one, given a single routine declaration")
  }

  switch v := declaration.(type) {
  case *Ast.SubroutineDeclaration:
    subroutineDeclaration = declaration.(*Ast.SubroutineDeclaration)
    if subroutineDeclaration.Name != subroutineName {
      t.Error("Variable name doesn't match, got", subroutineDeclaration.Name, "expected", subroutineName)
    }
  default:
    t.Error("Node", v, "is not Ast.SubroutineDeclaration")
  }

}

func TestSimpleSubroutineDeclaration( t *testing.T ) {
  source := `
  subrutina a()
  inicio
  fin
  `
  f, err := parser.Parse(source)

  testSubroutineDeclaration(f, t, "a")

  if err != nil {
    panic(err)
  }
}
