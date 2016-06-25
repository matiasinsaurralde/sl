package parser_test

import(
  "github.com/matiasinsaurralde/sl/parser"
  "github.com/matiasinsaurralde/sl/ast"

  "testing"
)

func testGenericDeclaration( f *Ast.File, t *testing.T, variableName string, variableType string, value string ) {

  var genericDeclaration *Ast.GenericDeclaration
  declaration := *(&f.Nodes[0])

  if len(f.Nodes) == 0 {
    t.Error("Ast.File Nodes is empty")
  }

  if len(f.Nodes) > 1 {
    t.Error("Ast.File Nodes is larger than one, given a single declaration")
  }

  switch v := declaration.(type) {
  case *Ast.GenericDeclaration:
    genericDeclaration = declaration.(*Ast.GenericDeclaration)
    if genericDeclaration.Name != variableName {
      t.Error("Variable name doesn't match, got", genericDeclaration.Name, "expected", variableName)
    }
  default:
    t.Error("Node", v, "is not Ast.GenericDeclaration")
  }

}

func TestInlineSingleDeclaration( t *testing.T ) {
  source := `var n:numerico`
  f, err := parser.Parse(source)

  testGenericDeclaration(f, t, "n", "numerico", "")

  if err != nil {
    panic(err)
  }
}

func TestMultilineSingleDeclaration( t *testing.T ) {
  source := `var
              n:numerico`
  f, err := parser.Parse(source)

  testGenericDeclaration(f, t, "n", "numerico", "")

  if err != nil {
    panic(err)
  }
}

func TestInlineSingleDeclarationWithIntValue( t *testing.T ) {
  source := `var n = 1`
  f, err := parser.Parse(source)

  testGenericDeclaration(f, t, "n", "", "1")

  if err != nil {
    panic(err)
  }
}

func TestInlineSingleDeclarationWithStringValue( t *testing.T ) {
  source := `var n = "test"`
  f, err := parser.Parse(source)

  declaration := *(&f.Nodes[0])

  switch v := declaration.(type) {
  case *Ast.GenericDeclaration:
  default:
    t.Error("Node", v, "is not Ast.GenericDeclaration")
  }
  if err != nil {
    panic(err)
  }
}

func TestMultilineDeclarationWithIntValue( t *testing.T ) {
  source := `var
              a = 0
              b = 10`
  _, err := parser.Parse(source)
  if err != nil {
    panic(err)
  }
}

func TestMultilineDeclarationWithStringValue( t *testing.T ) {
  source := `var
              a = "a"
              b = "b"`
  _, err := parser.Parse(source)
  if err != nil {
    panic(err)
  }
}

func TestMultilineDeclarationWithMixedValues( t *testing.T ) {
  source := `var
              a = "a"
              b = 100`
  _, err := parser.Parse(source)
  if err != nil {
    panic(err)
  }
}
